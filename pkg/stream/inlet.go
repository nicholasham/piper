package stream

import (
	"sync"
)

type Inlet struct {
	name               string
	in                 chan Element
	done               chan struct{}
	once               sync.Once
	completionSignaled bool
}

func (i *Inlet) In() <-chan Element {
	return i.in
}

func (i *Inlet) Complete() {
	i.once.Do(func() {
		close(i.done)
		i.completionSignaled = true
	})
}

func (i *Inlet) CompletionSignaled() bool {
	return i.completionSignaled
}

func (i *Inlet) WireTo(outlet *Outlet) *Inlet {
	i.in = outlet.out
	i.done = outlet.done
	return i
}

func NewInlet(attributes *StageAttributes) *Inlet {
	return &Inlet{
		name: attributes.Name + ".in",
		in:   make(chan Element),
		done: make(chan struct{}),
	}
}
