package stream

import (
	"fmt"
	"sync"
)

type Stream interface {
	Reader() Reader
	Writer() Writer
}

// verify stream implements Stream interface
var _ Stream = (*stream)(nil)
var _ Writer = (*stream)(nil)
var _ Reader = (*stream)(nil)

type stream struct {
	name string
	elements           chan Element
	done               chan struct{}
	completeOnce       sync.Once
	closeOnce          sync.Once
	completionSignaled bool
	closed             bool
	sync.Mutex
}

func (s *stream) Done() chan struct{} {
	return s.done
}

func (s *stream) Closed() bool {
	return s.closed
}

func (s *stream) Completing() bool {
	return s.completionSignaled
}

func (s *stream) Read() <-chan Element {
	return s.elements
}

func (s *stream) Complete() {
	s.completeOnce.Do(func() {
		fmt.Println("Completing stream " + s.name)
		close(s.done)
		s.completionSignaled = true
	})
}

func (s *stream) Close() {
	s.closeOnce.Do(func() {
		fmt.Println("Done stream  "+ s.name)
		close(s.elements)
		s.closed = true
	})
}

func (s *stream) Send(element Element) {
	s.Lock()
	defer s.Unlock()
	s.elements <- element
}

func (s *stream) Reader() Reader {
	return s
}

func (s *stream) Writer() Writer {
	return s
}

func NewStreamOld(name string) Stream {
	return &stream{
		name: name,
		done:     make(chan struct{}),
		elements: make(chan Element),
	}
}
