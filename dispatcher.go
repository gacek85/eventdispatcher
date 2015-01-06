// Package eventdispatcher contains a set of tools making up a simple and
// reliable event dispatcher
package eventdispatcher

import (
	"reflect"
	"strings"
	"sync"
)

const (
	DefaultDispatcherKey = "event_dispatcher"
)

// Listener type for defining functions as listeners
type Listener func(Event)

// Dispatcher interface defines the event dispatcher behavior
type Dispatcher interface {

	// Dispatch dispatches the event and returns it after all listeners do
	// their jobs.
	Dispatch(e Event) Event

	// On registers a listener for given event name.
	On(n string, l Listener)

	// Once registers a listener to be executed only once. The first param
	// n is the name of the event the listener will listen on, second is
	// the Listener type function.
	Once(n string, l Listener)

	// Off removes the registered event listener for given event name.
	Off(n string, l Listener)

	// RemoveAll removes all listeners for given name.
	OffAll(n string)

	// HasListeners returns true if any listener for given event name has
	// been assigned and false otherwise. This applies also to once triggered
	// listeners registered with `One` method
	HasListeners(n string) bool
}

type listenersCollection []Listener

// The EventDispatcher type is the default implementation of the
// DispatcherInterface
type EventDispatcher struct {
	sync.RWMutex
	listeners map[string]listenersCollection
}

// Forces the instance to be aware of event dispatcher
type DispatcherAware interface {

	// Dispatcher provides the event dispatcher instance pointer
	Dispatcher() Dispatcher
}

// On registers a listener for given event name.
func (d *EventDispatcher) On(n string, l Listener) {
	names := getNames(n)
	for _, name := range names {
		on(d, name, l)
	}
}

// getNames splits the given n string with space and returns a slice of
// event names strings
func getNames(n string) []string {
	names := strings.Split(n, " ")
	var results []string
	for _, name := range names {
		if name != "" {
			results = append(results, name)
		}
	}

	return results
}

// on binds listener to given event name n
func on(d *EventDispatcher, n string, l Listener) {
	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()
	d.listeners[n] = append(d.listeners[n], l)
}

// Once registers a listener to be executed only once. The first param
// n is the name of the event the listener will listen on, second is
// the Listener type function.
func (d *EventDispatcher) Once(n string, l Listener) {
	names := getNames(n)
	for _, name := range names {
		nl := executeRemove(d, name, l) // Create a new listener that removes given listener after calling it
		on(d, n, nl)
	}
}

func executeRemove(d *EventDispatcher, n string, l Listener) Listener {
	var nl func(e Event)
	nl = func(e Event) {
		l(e)
		d.RWMutex.RUnlock() // The dispatcher is locked in the Dispatch method, need to unlock it
		d.Off(n, nl)
		d.RWMutex.RLock()
	}

	return nl
}

// Off removes the registered event listener for given event name.
func (d *EventDispatcher) Off(n string, l Listener) {
	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()

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
func (d *EventDispatcher) OffAll(n string) {
	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()

	_, ok := d.listeners[n]
	if ok != false {
		delete(d.listeners, n)
	}
}

// HasListeners returns true if any listener for given event name has
// been assigned and false otherwise. This applies also to once triggered
// listeners registered with `One` method
func (d *EventDispatcher) HasListeners(n string) bool {
	listeners, ok := d.listeners[n]
	if ok == false {
		return false
	}

	return len(listeners) != 0
}

// Dispatch dispatches the event and returns it after all listeners do their jobs
func (d *EventDispatcher) Dispatch(e Event) Event {
	d.RWMutex.RLock()
	defer d.RWMutex.RUnlock()

	return dispatch(d, e)
}

// dispatch takes all registered listeners for given event name
// and dispatches the event
func dispatch(d *EventDispatcher, e Event) Event {
	for _, l := range d.listeners[e.Name()] {
		l(e)
	}

	return e
}

// Inner registry of event dispatcher instances
var dispatchers map[string]*EventDispatcher

// GetDispatcher provides event dispatcher for given key string. If the
// key string is nil, takes the default key
func GetDispatcher(k interface{}) *EventDispatcher {
	var key string
	if dispatchers == nil {
		dispatchers = make(map[string]*EventDispatcher)
	}
	if k == nil {
		key = DefaultDispatcherKey
	} else {
		key = k.(string)
	}

	return getDispatcher(key)
}

func getDispatcher(k string) *EventDispatcher {
	d, ok := dispatchers[k]
	if ok == false {
		dispatchers[k] = NewDispatcher()
		d = dispatchers[k]
	}
	return d
}

// NewDispatcher creates a new instance of event dispatcher
func NewDispatcher() *EventDispatcher {
	return &EventDispatcher{
		listeners: make(map[string]listenersCollection),
	}
}
