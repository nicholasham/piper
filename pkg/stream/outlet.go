package stream

import (
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

func (o *Outlet) send(element Element) {
	o.Lock()
	defer o.Unlock()
	o.out <- element
}

func (o *Outlet) SendValue(value interface{}) {
	o.send(Value(value))
}

func (o *Outlet) SendError(err error) {
	o.send(Error(err))
}

func (o *Outlet) Done() chan struct{} {
	return o.done
}

func (o *Outlet) Close() {
	o.once.Do(func() {
		o.closed = true
		close(o.out)
	})
}

func NewOutlet(attributes *StageAttributes) *Outlet {
	return &Outlet{
		name: attributes.Name + ".out",
		out:  createChannel(attributes),
		done: make(chan struct{}),
	}
}

func createChannel(attributes *StageAttributes) chan Element {
	if attributes.OutputBufferSize > 0 {
		return make(chan Element, attributes.OutputBufferSize)
	}
	return make(chan Element)
}
