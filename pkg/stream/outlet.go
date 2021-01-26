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

func NewOutlet(state *StageState) *Outlet {
	return &Outlet{
		name: state.Name + ".outputStage",
		out:  createChannel(state),
		done: make(chan struct{}),
	}
}

func createChannel(options *StageState) chan Element {
	if options.OutputBufferSize > 0 {
		return make(chan Element, options.OutputBufferSize)
	}
	return make(chan Element)
}
