package pipe

import (
	"github.com/nicholasham/piper/pkg/core"
)


type Stream struct {
	name     string
	elements chan core.Result
	done     chan struct{} // signal channel
}

func (s *Stream) Done() {
	if !s.IsDone() {
		close(s.done)
	}
}

func (s *Stream) Close()  {
	s.Done()
	close(s.elements)
}

func (s *Stream) Receiver() *Receiver {
	return & Receiver{stream: s}
}

func (s *Stream) Sender() *Sender {
	return & Sender{stream: s}
}

func (s *Stream) IsDone() bool {
	select {
	case <-s.done:
		return true
	default:
	}
	return false
}

type Receiver struct {
	stream *Stream
}

func (r *Receiver) Done() {
	r.stream.Done()
}

func (r *Receiver) Receive() chan core.Result {
	return r.stream.elements
}

type Sender struct {
	stream *Stream
}

func (s *Sender) Close() {
	s.stream.Close()
}

func (s *Sender) Done() chan struct{} {
	return s.stream.done
}

func (s *Sender) IsDone() bool {
	return s.stream.IsDone()
}

func (s *Sender) TrySend(element core.Result) bool {
	select {
	case <- s.stream.done:
		return false
	case s.stream.elements <- element:
		return true
	}
}

func NewStream(name string) *Stream {
	return &Stream{
		name :    name,
		elements: make(chan core.Result),
		done:     make(chan struct{}),
	}
}
