// Package eventdispatcher contains a set of tools making up a simple and
// reliable event dispatcher
package eventdispatcher

// Event is an interface used by event dispatcher. Contains name and more custom data
// May be forced to stop being propagated
type Event interface {

	// Returns the event name
	Name() string

	// IsPropagationStopped informs weather the event should
	// be further propagated or not
	IsPropagationStopped() bool

	// StopPropagation makes the event no longer
	// propagate.
	StopPropagation()
}

// ParamsEvent is the default implementation of Event interface. Contains additional
// string parameters
type ParamsEvent struct {
	name                 string
	isPropagationStopped bool
	params               map[string]interface{}
}

// Name returns the name of the event
func (event *ParamsEvent) Name() string {
	return event.name
}

// IsPropagationStopped informs weather the event should
// be further propagated or not
func (event *ParamsEvent) IsPropagationStopped() bool {
	return event.isPropagationStopped
}

// StopPropagation sets a flag that make the event no longer
// propagate.
func (event *ParamsEvent) StopPropagation() {
	event.isPropagationStopped = true
}

// AddParam registers a parameter for the event.
// Returns this event instance
func (event *ParamsEvent) SetParam(k string, v interface{}) *ParamsEvent {
	event.params[k] = v
	return event
}

// RemoveParam deletes a param with given key. Does nothing, if
// the param does not exst. Returns this event instance
func (event *ParamsEvent) RemoveParam(k string) *ParamsEvent {
	if event.HasParam(k) {
		delete(event.params, k)
	}
	return event
}

// HasParam defines if a param with given key exists. Returns a boolean value
func (event *ParamsEvent) HasParam(k string) bool {
	_, ok := event.params[k]
	return ok
}

// GetParam returns a parameter value for given key. If the param does not exist,
// returns an empty string. Second value returned contains boolean value that says
// if the param existed.
func (event *ParamsEvent) GetParam(k string) (value interface{}, ok bool) {
	v, ok := event.params[k]
	if ok == false {
		return "", false
	}
	return v, ok
}

// NewParamsEvent is a factory for creating a basic event
func NewParamsEvent(n string) *ParamsEvent {
	p := make(map[string]interface{})
	e := ParamsEvent{n, false, p} // Propagation never stopped by default
	return &e
}
