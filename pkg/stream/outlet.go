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

func (o *Outlet) Send(element Element) {
	o.Lock()
	defer o.Unlock()
	o.out <- element
}

func (o *Outlet) SendValue(value interface{}) {
	o.Send(Value(value))
}

func (o *Outlet) SendError(err error) {
	o.Send(Error(err))
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

func NewOutlet(options *StageOptions) *Outlet {
	return &Outlet{
		name: options.Name + ".out",
		out:  createChannel(options),
		done: make(chan struct{}),
	}
}

func createChannel(options *StageOptions) chan Element {
	if options.OutputBufferSize > 0 {
		return make(chan Element, options.OutputBufferSize)
	}
	return make(chan Element)
}
