package piper

import "sync"

type concurrentBuffer struct {
	sync.RWMutex
	items []interface{}
}

func (b *concurrentBuffer) Append(item interface{}) {
	b.Lock()
	defer b.Unlock()
	b.items = append(b.items, item)
}

func (b *concurrentBuffer) Clear() {
	b.Lock()
	defer b.Unlock()
	b.items = []interface{}{}
	return
}

func (b *concurrentBuffer) Result() []interface{} {
	return b.items
}

func (b *concurrentBuffer) Count() int {
	return len(b.items)
}

type ConcurrentBuffer interface {
	Append(item interface{})
	Clear()
	Result() []interface{}
	Count() int
}

func NewConcurrentBuffer() ConcurrentBuffer {
	return &concurrentBuffer{}
}
