package events

import "sync/atomic"

type Event struct {
	Type     string      `json:"type"`
	ServerID string      `json:"serverId"`
	Data     interface{} `json:"data,omitempty"`
}

type Subscriber struct {
	id int64
	Ch chan Event
}

type Bus struct {
	subs   map[int64]*Subscriber
	nextID int64
	mu     chan func()
}

func NewBus() *Bus {
	b := &Bus{
		subs: make(map[int64]*Subscriber),
		mu:   make(chan func(), 64),
	}
	go func() {
		for f := range b.mu {
			f()
		}
	}()
	return b
}

func (b *Bus) do(f func()) { b.mu <- f }

func (b *Bus) Subscribe() *Subscriber {
	sub := &Subscriber{Ch: make(chan Event, 256)}
	b.do(func() {
		id := atomic.AddInt64(&b.nextID, 1)
		sub.id = id
		b.subs[id] = sub
	})
	return sub
}

func (b *Bus) Unsubscribe(s *Subscriber) {
	b.do(func() {
		delete(b.subs, s.id)
		close(s.Ch)
	})
}

func (b *Bus) Publish(ev Event) {
	b.do(func() {
		for _, s := range b.subs {
			select {
			case s.Ch <- ev:
			default:
				// drop event if subscriber is slow
			}
		}
	})
}
