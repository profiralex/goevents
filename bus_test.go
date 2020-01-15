package goevents_test

import "testing"
import "github.com/profiralex/goevents"
import "github.com/stretchr/testify/mock"

import "time"

const (
	testEvent1 = "testEvent1"
	testEvent2 = "testEvent2"
)

func TestNotify(t *testing.T) {
	event := int64(1)

	listener := new(testListener)
	listener.On("Notify", event)

	bus := goevents.NewBus(100)
	bus.Subscribe(testEvent1, listener)
	bus.Notify(testEvent1, event)

	time.Sleep(10 * time.Millisecond)

	listener.AssertExpectations(t)
}

func TestNotify2(t *testing.T) {
	event := int64(1)
	complexEvent := map[string]string{"field1": "value1", "field2": "value2"}

	listener1 := new(testListener)
	listener1.On("Notify", event)

	listener2 := new(testListener)
	listener2.On("Notify", event)

	listener3 := new(testListener)
	listener3.On("Notify", complexEvent)

	bus := goevents.NewBus(100)
	bus.Subscribe(testEvent1, listener1)
	bus.Subscribe(testEvent1, listener2)
	bus.Subscribe(testEvent2, listener3)

	bus.Notify(testEvent1, event)
	bus.Notify(testEvent1, event)
	bus.Notify(testEvent2, complexEvent)

	time.Sleep(100 * time.Millisecond)

	listener1.AssertExpectations(t)
	listener1.AssertNumberOfCalls(t, "Notify", 2)
	listener2.AssertExpectations(t)
	listener2.AssertNumberOfCalls(t, "Notify", 2)
	listener3.AssertExpectations(t)
	listener3.AssertNumberOfCalls(t, "Notify", 1)
}

type testListener struct {
	mock.Mock
}

func (o *testListener) Notify(event interface{}) {
	o.Called(event)
}
