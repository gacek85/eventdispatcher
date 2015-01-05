// Package eventdispatcher contains a set of tools making up a simple and
// reliable event dispatcher
package eventdispatcher

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Tests the name getter
func TestName(t *testing.T) {
	assert := assert.New(t)
	e := getTestEvent()
	n := getTestEventName()
	assert.Equal(n, e.Name(), fmt.Sprintf("The event name provided %q is different than the one taken from the event: %q", n, e.Name()))
}

func getTestEvent() ParamsEvent {
	e := NewParamsEvent(getTestEventName())
	return *e
}

func getTestEventName() string {
	n := "test_name"
	return n
}

func TestStopPropagation(t *testing.T) {
	assert := assert.New(t)
	e := getTestEvent()
	assert.False(e.IsPropagationStopped(), fmt.Sprintf("The event propagation is stopped before calling %q", "IsPropagationStopped"))
	e.StopPropagation()
	assert.True(e.IsPropagationStopped(), fmt.Sprintf("The event propagation is not stopped after calling %q", "IsPropagationStopped"))
}

type testType struct {
	foo string
	bar int64
}

func TestParams(t *testing.T) {
	assert := assert.New(t)
	e := getTestEvent()
	p := "test value"
	k := "test_key"

	// Has nonexisting param
	assert.False(e.HasParam(k), fmt.Sprintf("The event does not contain the param %s", k))

	// Get nonexisting param
	v, ok := e.GetParam(k)
	assert.Equal("", v, fmt.Sprintf("Expected value for nonexisting key is `%s`, returned value %s", "", v))
	assert.False(ok, fmt.Sprintf("%s expects second returned value to be false if param does not exists.", "GetParam"))

	re := e.SetParam(k, p)
	assert.Equal(&e, re, fmt.Sprintf("The %s method should return same event instance for chaining!", "SetParam"))

	// Has existing param
	assert.True(e.HasParam(k), fmt.Sprintf("The event should contain the param %s with value %s", k, p))

	// Get existing param
	v, ok = e.GetParam(k)
	assert.Equal(p, v.(string), fmt.Sprintf("Expected value %s, returned value %s", p, v.(string)))
	assert.True(ok, fmt.Sprintf("%s expects second returned value to be true if param exists.", "GetParam"))

	// Remove param
	re = e.RemoveParam(k)
	assert.Equal(&e, re, fmt.Sprintf("The %s method should return same event instance for chaining!", "SetParam"))

	// Again check nonexisting params
	assert.False(e.HasParam(k), fmt.Sprintf("The event does not contain the param %s", k))

	// Testing setting parameter with non-string value
	ns := testType{"test", 1}
	e.SetParam(k, ns)
	rns, _ := e.GetParam(k)
	assert.Equal(ns, rns.(testType), fmt.Sprintf("The event does not contain valid param %s", k))
}
