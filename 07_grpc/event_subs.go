package main

import "sync"

func NewEventSubs() *EventSubs {
	return &EventSubs{
		Subs: make(map[interface{}]chan *Event, 32),
	}
}

type EventSubs struct {
	Subs map[interface{}]chan *Event
	mu   sync.RWMutex
}

func (es *EventSubs) Subscribe(subscriber interface{}) chan *Event {
	ch := make(chan *Event)

	es.mu.Lock()
	defer es.mu.Unlock()
	es.Subs[subscriber] = ch

	return ch
}

func (es *EventSubs) Unsubscribe(subscriber interface{}) {
	es.mu.Lock()
	defer es.mu.Unlock()

	close(es.Subs[subscriber])
	delete(es.Subs, subscriber)
}

func (es *EventSubs) Publish(event *Event) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	for _, ch := range es.Subs {
		ch <- event
	}
}
