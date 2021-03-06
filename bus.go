package goevents

import "sync"

// Listener interface to listen for events
type Listener interface {
	Notify(interface{})
}

type event struct {
	name string
	data interface{}
}

// Bus manages the listeners and events
type Bus struct {
	listeners map[string]([]Listener)
	mux       sync.Mutex
	channel   chan event
}

// NewBus creates and initializes a new bus
func NewBus(cap int) *Bus {
	b := Bus{}
	b.channel = make(chan event, cap)
	b.listeners = make(map[string]([]Listener))

	go b.listenToEvents()

	return &b
}

// Subscribe listener to event
func (b *Bus) Subscribe(event string, listener Listener) {
	b.mux.Lock()
	found := false
	for _, l := range b.listeners[event] {
		if l == listener {
			found = true
			break
		}
	}

	if !found {
		b.listeners[event] = append(b.listeners[event], listener)
	}
	b.mux.Unlock()
}

// Unsubscribe listener from events
func (b *Bus) Unsubscribe(event string, listener Listener) {
	b.mux.Lock()
	listeners, ok := b.listeners[event]
	if ok {
		newListeners := []Listener{}
		for _, l := range listeners {
			if l != listener {
				newListeners = append(newListeners, l)
			}
		}
		b.listeners[event] = newListeners
	}
	b.mux.Unlock()
}

// Notify listeners about event
func (b *Bus) Notify(eventName string, data interface{}) {
	b.channel <- event{name: eventName, data: data}
}

func (b *Bus) listenToEvents() {
	for event := range b.channel {
		listeners := []Listener{}

		b.mux.Lock()
		for _, listener := range b.listeners[event.name] {
			listeners = append(listeners, listener)
		}
		b.mux.Unlock()

		for _, listener := range listeners {
			listener.Notify(event.data)
		}
	}
}
