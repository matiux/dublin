package eventdispatcher

// EventDispatcher is a base type fot an event dispatcher
type EventDispatcher interface {

	Dispatch(eventName string, arguments map[string]interface{})

	AddListener(eventName string, callable func(args ...interface{}))
}
