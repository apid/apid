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
	em.dispatchers[selector].Send(event)
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
		d := &dispatcher{sync.RWMutex{}, em, selector, nil, nil}
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
	handler := &funcWrapper{handlerFunc}
	em.Listen(selector, handler)
}

func (em *eventManager) Close() {
	log.Debugf("Closing")
	dispatchers := em.dispatchers
	em.dispatchers = nil
	for _, dispatcher := range dispatchers {
		dispatcher.Close()
	}
}

type dispatcher struct {
	sync.RWMutex
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

func (d *dispatcher) Send(e apid.Event) {
	if d != nil {
		d.channel <- e
	}
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
					if d.selector != apid.EventDeliveredSelector &&
						d.em.HasListeners(apid.EventDeliveredSelector) {

						e := apid.EventDeliveryEvent{
							Description: "event complete",
							Selector:    d.selector,
							Event:       event,
						}
						d.em.Emit(apid.EventDeliveredSelector, e)
					}
					log.Debugf("delivery complete")
				}
			}

		}
	}()
}

type funcWrapper struct {
	f apid.EventHandlerFunc
}

func (r *funcWrapper) Handle(e apid.Event) {
	r.f(e)
}
