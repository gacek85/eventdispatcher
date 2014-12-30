// Package eventdispatcher contains a set of tools making up a simple and
// reliable event dispatcher
package eventdispatcher

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

const (
	TestEventName = "test_event"
)


func TestOnOff (t *testing.T) {
	assert := assert.New(t)
	d := NewDispatcher()
	assert.False(d.HasListeners(TestEventName), fmt.Sprintf("No listeners assigned yet for %s!", TestEventName))
	l := func (e Event) {
		fmt.Sprintf("Event name: %s", e.Name())
	}
	d.On(TestEventName, l)
	assert.True(d.HasListeners(TestEventName), fmt.Sprintf("There should be listeners assigned for %s!", TestEventName))
	d.Off(TestEventName, l)
	assert.False(d.HasListeners(TestEventName), fmt.Sprintf("No listeners assigned for %s!", TestEventName))
}

func TestOnce (t *testing.T) {
	assert := assert.New(t)
	d := NewDispatcher()
	d.Once(TestEventName, func (e Event) {
		fmt.Sprintf("Event name: %s", e.Name())
	})
	assert.True(d.HasListeners(TestEventName), fmt.Sprintf("There should be one listener assigned for one call for %s!", TestEventName))
	e := NewParamsEvent(TestEventName)
	d.Dispatch(e)
	assert.False(d.HasListeners(TestEventName), fmt.Sprintf("The listener called should unbind itself for %s!", TestEventName))
}


func TestDispatch (t *testing.T) {
	assert := assert.New(t)
	d := NewDispatcher()
	var c int
	c = 0
	for i := 1;  i<=5; i++ {
		d.On(TestEventName, func (e Event) {
			c++
		})
	}
	e := NewParamsEvent(TestEventName)
	d.Dispatch(e)
	assert.Equal(5, c, "Invalid listeners calls number!")
}


func TestOffAll (t *testing.T) {
	assert := assert.New(t)
	d := NewDispatcher()
	var c int
	c = 0
	for i := 1;  i<=5; i++ {
		d.On(TestEventName, func (e Event) {
			c++
		})
	}
	assert.True(d.HasListeners(TestEventName), fmt.Sprintf("There should be listeners assigned for %s!", TestEventName))
	d.OffAll(TestEventName)
	assert.False(d.HasListeners(TestEventName), fmt.Sprintf("All event listeners for %s should be removed!", TestEventName))
}


func TestGetDispatcher (t *testing.T) {
	assert := assert.New(t)
	d := GetDispatcher(nil)
	dn := GetDispatcher(nil)
	assert.Equal(d, dn, "The event dispatchers should be the same instance pointers!")
	_ = GetDispatcher("foo")
}
