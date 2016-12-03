package events

import (
	"github.com/30x/apid"
	"sync"
)

// events published to a given channel are processed entirely in order, though delivery to listeners is async

type eventManager struct {
	dispatchers map[apid.EventSelector]*dispatcher
}

func (em *eventManager) Emit(selector apid.EventSelector, event apid.Event) {
	log.Debugf("emit selector: '%s' event: %s", selector, event)
	if !em.dispatchers[selector].Send(event) {
		em.sendDelivered(selector, event, 0) // in case of no dispatcher
	}
}

func (em *eventManager) EmitWithCallback(selector apid.EventSelector, event apid.Event, callback apid.EventHandlerFunc) {
	log.Debugf("emit with callback selector: '%s' event: %s", selector, event)

	handler := &funcWrapper{em, nil}
	handler.HandlerFunc = func(e apid.Event) {
		if ede, ok := e.(apid.EventDeliveryEvent); ok {
			if ede.Event == event {
				em.StopListening(apid.EventDeliveredSelector, handler)
				callback(e)
			}
		}
	}

	em.Listen(apid.EventDeliveredSelector, handler)
	em.Emit(selector, event)
}

func (em *eventManager) HasListeners(selector apid.EventSelector) bool {
	return em.dispatchers[selector].HasHandlers()
}

func (em *eventManager) Listen(selector apid.EventSelector, handler apid.EventHandler) {
	log.Debugf("listen: '%s' handler: %s", selector, handler)
	if em.dispatchers == nil {
		em.dispatchers = make(map[apid.EventSelector]*dispatcher)
	}
	list := em.dispatchers[selector]
	if list == nil {
		d := &dispatcher{sync.Mutex{}, em, selector, nil, nil}
		em.dispatchers[selector] = d
	}
	em.dispatchers[selector].Add(handler)
}

func (em *eventManager) StopListening(selector apid.EventSelector, handler apid.EventHandler) {
	log.Debugf("stop listening: '%s' handler: %s", selector, handler)
	if em.dispatchers == nil {
		return
	}
	em.dispatchers[selector].Remove(handler)
}

func (em *eventManager) ListenFunc(selector apid.EventSelector, handlerFunc apid.EventHandlerFunc) {
	log.Debugf("listenFunc: '%s' handler: %s", selector, handlerFunc)
	handler := &funcWrapper{em, handlerFunc}
	em.Listen(selector, handler)
}

func (em *eventManager) ListenOnceFunc(selector apid.EventSelector, handlerFunc apid.EventHandlerFunc) {
	log.Debugf("listenOnceFunc: '%s' handler: %s", selector, handlerFunc)
	handler := &funcWrapper{em, nil}
	handler.HandlerFunc = func(event apid.Event) {
		em.StopListening(selector, handler)
		handlerFunc(event)
	}
	em.Listen(selector, handler)
}

func (em *eventManager) Close() {
	log.Debugf("Closing %d dispatchers", len(em.dispatchers))
	dispatchers := em.dispatchers
	em.dispatchers = nil
	for _, dispatcher := range dispatchers {
		dispatcher.Close()
	}
}

func (em *eventManager) sendDelivered(selector apid.EventSelector, event apid.Event, count int) {
	if selector != apid.EventDeliveredSelector {
		ede := apid.EventDeliveryEvent{
			Description: "event complete",
			Selector:    selector,
			Event:       event,
			Count:       count,
		}
		em.dispatchers[apid.EventDeliveredSelector].Send(ede)
	}
}

type dispatcher struct {
	sync.Mutex
	em       *eventManager
	selector apid.EventSelector
	channel  chan apid.Event
	handlers []apid.EventHandler
}

func (d *dispatcher) Add(h apid.EventHandler) {
	d.Lock()
	defer d.Unlock()
	if d.handlers == nil {
		d.handlers = []apid.EventHandler{h}
		d.channel = make(chan apid.Event, config.GetInt(configChannelBufferSize))
		d.startDelivery()
		return
	}
	cp := make([]apid.EventHandler, len(d.handlers)+1)
	copy(cp, d.handlers)
	cp[len(d.handlers)] = h
	d.handlers = cp
}

func (d *dispatcher) Remove(h apid.EventHandler) {
	d.Lock()
	defer d.Unlock()
	for i := len(d.handlers) - 1; i >= 0; i-- {
		ih := d.handlers[i]
		if h == ih {
			d.handlers = append(d.handlers[:i], d.handlers[i+1:]...)
			return
		}
	}
}

func (d *dispatcher) Close() {
	close(d.channel)
}

func (d *dispatcher) Send(e apid.Event) bool {
	if d != nil {
		d.channel <- e
		return true
	}
	return false
}

func (d *dispatcher) HasHandlers() bool {
	return d != nil && len(d.handlers) > 0
}

func (d *dispatcher) startDelivery() {
	go func() {
		for {
			select {
			case event := <-d.channel:
				if event != nil {
					log.Debugf("delivering %v to %v", event, d.handlers)
					if len(d.handlers) > 0 {
						var wg sync.WaitGroup
						for _, h := range d.handlers {
							handler := h
							wg.Add(1)
							go func() {
								defer wg.Done()
								handler.Handle(event) // todo: recover on error?
							}()
						}
						log.Debugf("waiting for handlers")
						wg.Wait()
					}
					d.em.sendDelivered(d.selector, event, len(d.handlers))
					log.Debugf("delivery complete")
				}
			}

		}
	}()
}

type funcWrapper struct {
	*eventManager
	HandlerFunc apid.EventHandlerFunc
}

func (r *funcWrapper) Handle(e apid.Event) {
	r.HandlerFunc(e)
}
