package stream

import "sync"

type Stream interface {
	Reader() Reader
	Writer() Writer
}

type Reader interface {
	Elements() <-chan Element
	Complete()
	Completing() bool
}

type Writer interface {
	Close()
	SendValue(value interface{})
	SendError(err error)
	Closed() bool
	Done() chan struct{}
}

// verify stream implements Stream interface
var _ Stream = (*stream)(nil)
var _ Writer = (*stream)(nil)
var _ Reader = (*stream)(nil)

type stream struct {
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

func (s *stream) Elements() <-chan Element {
	return s.elements
}

func (s *stream) Complete() {
	s.completeOnce.Do(func() {
		close(s.done)
		s.completionSignaled = true
	})
}

func (s *stream) Close() {
	s.closeOnce.Do(func() {
		s.closed = true
		close(s.elements)
	})
}

func (s *stream) SendValue(value interface{}) {
	s.send(Value(value))
}

func (s *stream) SendError(err error) {
	s.send(Error(err))
}

func (s *stream) send(element Element) {
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

func NewStream() Stream {
	return &stream{
		done:     make(chan struct{}),
		elements: make(chan Element),
	}
}
