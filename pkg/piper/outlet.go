package piper

import (
	"github.com/nicholasham/piper/pkg/piper/attribute"
	"sync"

)

type Outlet struct {
	name   string
	out    chan Element
	done   chan struct{}
	closed bool
	once   sync.Once
	sync.Mutex
}

func (receiver *Outlet) Out() chan Element {
	return receiver.out
}

func (receiver *Outlet) Send(element Element) {
	receiver.Lock()
	defer receiver.Unlock()
	receiver.out <- element
}

func (receiver *Outlet) Done() chan struct{} {
	return receiver.done
}

func (receiver *Outlet) Close() {
	receiver.once.Do(func() {
		receiver.closed = true
		close(receiver.out)
	})
}

func NewOutletOld(name string, options ...Option) *Outlet {
	return &Outlet{
		name: name + ".out",
		out:  CreateChannel(options),
		done: make(chan struct{}),
	}
}

func NewOutlet(stageAttributes *attribute.StageAttributes) *Outlet {
	return &Outlet{
		name: stageAttributes.Name + ".out",
		out:  createChannel(stageAttributes),
		done: make(chan struct{}),
	}
}

func createChannel(stageAttributes *attribute.StageAttributes) chan Element {
	if stageAttributes.OutputBufferSize > 0 {
		return make(chan Element, stageAttributes.OutputBufferSize)
	}
	return make(chan Element)
}
