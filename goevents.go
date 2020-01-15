package goevents

import "sync"

const defaultCapacity = 100

var bus *Bus
var mux sync.Mutex

// Subscribe listener to event
func Subscribe(event string, listener Listener) {
	initBusIfNoBus()
	bus.Subscribe(event, listener)
}

// Unsubscribe listener from events
func Unsubscribe(event string, listener Listener) {
	initBusIfNoBus()
	bus.Unsubscribe(event, listener)
}

// Notify listeners about event
func Notify(event string, data interface{}) {
	initBusIfNoBus()
	bus.Notify(event, data)
}

func initBusIfNoBus() {
	if bus == nil {
		mux.Lock()
		if bus == nil {
			bus = NewBus(defaultCapacity)
		}
		mux.Unlock()
	}
}
