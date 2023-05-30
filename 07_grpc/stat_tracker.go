package main

import (
	"errors"
	"sync"
)

func NewStat() *Stat {
	return &Stat{
		ByMethod:   make(map[string]uint64, 32),
		ByConsumer: make(map[string]uint64, 32),
	}
}

func NewStatTracker() *StatTracker {
	return &StatTracker{
		Subs: make(map[interface{}]*Stat, 32),
	}
}

type StatTracker struct {
	Subs map[interface{}]*Stat
	mu   sync.Mutex
}

func (st *StatTracker) Subscribe(subscriber interface{}) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.Subs[subscriber] = NewStat()
}

func (st *StatTracker) Unsubscribe(subscriber interface{}) {
	st.mu.Lock()
	defer st.mu.Unlock()
	delete(st.Subs, subscriber)
}

func (st *StatTracker) Pull(subscriber interface{}) (*Stat, error) {
	st.mu.Lock()
	defer st.mu.Unlock()
	stat, ok := st.Subs[subscriber]
	if !ok {
		return nil, errors.New("subscriber does not exist")
	}
	st.Subs[subscriber] = NewStat()
	return stat, nil
}

func (st *StatTracker) Track(method string, consumer string) {
	st.mu.Lock()
	defer st.mu.Unlock()
	for _, s := range st.Subs {
		s.ByConsumer[consumer]++
		s.ByMethod[method]++
	}
}
