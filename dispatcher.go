// Package eventdispatcher contains a set of tools making up a simple and
// reliable event dispatcher
package eventdispatcher

import (
	"reflect"
	"sync"
)


const (
	DefaultDispatcherKey = "event_dispatcher"
)


// Dispatcher interface defines the event dispatcher behavior
type Dispatcher interface {
	
	// Dispatch dispatches the event and returns it after all listeners do
	// their jobs.
	Dispatch (e Event) Event
	
	
	// On registers a listener for given event name.
	On (n string, l Listener)
	
	
	// Once registers a listener to be executed only once. The first param
	// n is the name of the event the listener will listen on, second is
	// the Listener type function.
	Once (n string, l Listener)
	
	
	// Off removes the registered event listener for given event name.
	Off (n string, l Listener)
	
	
	// RemoveAll removes all listeners for given name.
	OffAll (n string)
	
	
	// HasListeners returns true if any listener for given event name has
	// been assigned and false otherwise. This applies also to once triggered
	// listeners registered with `One` method
	HasListeners (n string) bool
}


type listenersCollection []Listener


// The EventDispatcher type is the default implementation of the
// DispatcherInterface
type EventDispatcher struct {
	sync.RWMutex
	listeners map[string]listenersCollection
}


// On registers a listener for given event name.
func (d *EventDispatcher) On (n string, l Listener) {
	d.Lock()
	defer d.Unlock()
	d.listeners[n] = append(d.listeners[n], l)
}


// Once registers a listener to be executed only once. The first param
// n is the name of the event the listener will listen on, second is
// the Listener type function.
func (d *EventDispatcher) Once (n string, l Listener) {
	nl := executeRemove(d, n, l); 	// Create a new listener that removes
									// given listener after calling it
	d.On(n, nl)
}


func executeRemove (d *EventDispatcher, n string, l Listener) Listener {
	var nl func (e Event)
	nl = func (e Event) {
		l(e)
		d.RUnlock() // The dispatcher is locked in the Dispatch method, need to unlock it
		d.Off(n, nl)
	}
	
	return nl
}


// Off removes the registered event listener for given event name.
func (d *EventDispatcher) Off (n string, l Listener) {
	d.Lock()
	defer d.Unlock()
	
	p := reflect.ValueOf(l).Pointer()

	listeners := d.listeners[n]
	for i, l := range listeners {
		lp := reflect.ValueOf(l).Pointer()
		if lp == p {
			d.listeners[n] = append(listeners[:i], listeners[i+1:]...)
		}
	}
}


// RemoveAll removes all listeners for given name.
func (d *EventDispatcher) OffAll (n string) {
	d.Lock()
	defer d.Unlock()
	
	_, ok := d.listeners[n]
	if ok != false {
		delete(d.listeners, n)
	}
}


// HasListeners returns true if any listener for given event name has
// been assigned and false otherwise. This applies also to once triggered
// listeners registered with `One` method
func (d *EventDispatcher) HasListeners (n string) bool {
	listeners, ok := d.listeners[n]
	if ok == false {
		return false
	}
	
	return len(listeners) != 0
}


// Dispatch dispatches the event and returns it after all listeners do their jobs
func (d *EventDispatcher) Dispatch (e Event) Event {
	d.RLock()
	defer d.RUnlock()
	
	return doDispatch(d, e)
}


func doDispatch (d *EventDispatcher, e Event) Event {
	for _, l := range d.listeners[e.Name()] {
		l(e)
	}
	
	return e
}


// Inner registry of event dispatcher instances
var dispatchers map[string]*EventDispatcher



// Provides event dispatcher for given key string. If the 
// key string is nil, takes the default key
func GetDispatcher (k interface{}) *EventDispatcher {
	var key string
	if (dispatchers == nil) {
		dispatchers = make(map[string]*EventDispatcher)
	}
	if k == nil {
		key = DefaultDispatcherKey
	} else {
		key = k.(string)
	}
	
	return doGetDispatcher (key)
}


func doGetDispatcher (k string) *EventDispatcher {
	d, ok := dispatchers[k]
	if (ok == false) {
		dispatchers[k] = NewDispatcher()
		d = dispatchers[k]
	}
	return d
}


// NewDispatcher creates a new instance of event dispatcher
func NewDispatcher () *EventDispatcher {
	return &EventDispatcher{
		listeners : make(map[string]listenersCollection),
	}
}

