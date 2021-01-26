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

func (receiver *Inlet) In() <-chan Element {
	return receiver.in
}

func (receiver *Inlet) Complete() {
	receiver.once.Do(func() {
		close(receiver.done)
		receiver.completionSignaled = true
	})
}

func (receiver *Inlet) CompletionSignaled() bool {
	return receiver.completionSignaled
}

func (receiver *Inlet) WireTo(outlet *Outlet) *Inlet {
	receiver.in = outlet.out
	receiver.done = outlet.done
	return receiver
}

func NewInletOld(name string) *Inlet {
	return &Inlet{
		name: name + ".inputStage",
		in:   make(chan Element),
		done: make(chan struct{}),
	}
}

func NewInlet(stageAttributes *StageState) *Inlet {
	return &Inlet{
		name: stageAttributes.Name + ".inputStage",
		in:   make(chan Element),
		done: make(chan struct{}),
	}
}
