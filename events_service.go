package apid

type EventSelector string

type Event interface{}

type EventHandler interface {
	Handle(event Event)
}

type EventsService interface {
	// publish an event to the selector
	Emit(selector EventSelector, event Event)

	// when an event matching selector occurs, run the provided handler
	Listen(selector EventSelector, handler EventHandler)

	// shut it down
	Close()
}
