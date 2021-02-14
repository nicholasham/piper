package pipe

import (
	"fmt"
	"github.com/nicholasham/piper/pkg/core"
)


type Stream struct {
	name     string
	elements chan core.Result
	done     chan struct{} // signal channel
	closed   chan struct{}
}

func (s *Stream) Done() {
	select {
	case s.done <-struct{}{}:
		<-s.closed
	case <-s.closed:
	}
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

func (r *Receiver) Done(who string) {
	fmt.Println( fmt.Sprintf("%v has signaled that it is done with the %v receiver", who, r.stream.name))
	r.stream.Done()
}

func (r *Receiver) Receive() chan core.Result {
	return r.stream.elements
}

type Sender struct {
	stream *Stream
}

func (s *Sender) Close() {
	close(s.stream.closed)
	close(s.stream.elements)
}

func (s *Sender) Done() chan struct{} {
	return s.stream.done
}

func (s *Sender) IsDone() bool {
	return s.stream.IsDone()
}

func (s *Sender) Send(element core.Result) bool {
	select {
	case <- s.stream.done:
		fmt.Println( s.stream.name + " done...")
		return false
	case s.stream.elements <- element:
		return true
	}
}


func NewStream(name string) *Stream {
	return &Stream{
		name :    name,
		elements: make(chan core.Result, 10),
		closed:   make(chan struct{}),
		done:     make(chan struct{}),
	}
}

