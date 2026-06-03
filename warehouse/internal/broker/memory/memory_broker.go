package memory

import (
	"sync"
	"warehouse/internal/broker"
)

type MemoryBroker struct {
	mu     sync.RWMutex
	queues map[string][]chan []byte
}

func NewMemoryBroker() broker.MessageBroker {
	return &MemoryBroker{
		queues: make(map[string][]chan []byte),
	}
}

func (b *MemoryBroker) Publish(queue string, body []byte) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.queues[queue] {
		ch <- body
	}
	return nil
}

func (b *MemoryBroker) Subscribe(queue string, handler func([]byte) error) error {
	ch := make(chan []byte, 10)
	b.mu.Lock()
	b.queues[queue] = append(b.queues[queue], ch)
	b.mu.Unlock()

	go func() {
		for msg := range ch {
			if err := handler(msg); err != nil {
				// log error
			}
		}
	}()
	return nil
}
