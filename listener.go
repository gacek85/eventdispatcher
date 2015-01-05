// Package eventdispatcher contains a set of tools making up a simple and
// reliable event dispatcher
package eventdispatcher

// Listener type for defining functions as listeners
type Listener func(Event)
