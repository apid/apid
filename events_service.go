package apid

type EventSelector string

type Event interface{}

type EventHandler interface {
	Handle(event Event)
}

type EventHandlerFunc func(event Event)

type EventsService interface {
	// publish an event to the selector
	Emit(selector EventSelector, event Event)

	// publish an event to the selector, call the passed handler when all listeners have responded to the event
	EmitWithCallback(selector EventSelector, event Event, handler EventHandlerFunc)

	// when an event matching selector occurs, run the provided handler
	Listen(selector EventSelector, handler EventHandler)

	// when an event matching selector occurs, run the provided handler function
	ListenFunc(selector EventSelector, handler EventHandlerFunc)

	// when an event matching selector occurs, run the provided handler function and stop listening
	ListenOnceFunc(selector EventSelector, handler EventHandlerFunc)

	// remove a listener
	StopListening(selector EventSelector, handler EventHandler)

	// shut it down
	Close()
}

const EventDeliveredSelector EventSelector = "event delivered"

type EventDeliveryEvent struct {
	Description string
	Selector    EventSelector
	Event       Event
	Count       int
}

type PluginsInitializedEvent struct {
	Description string
	Plugins []PluginData
}

type PluginData struct {
	Name string
	Version string
	ExtraData map[string]interface{}
}
