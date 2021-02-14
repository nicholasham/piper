package stream

// https://go101.org/article/channel-closing.html

type Reader interface {
	Read() <-chan Element
	Complete()
	Completing() bool
}

type Writer interface {
	Send(element Element)
	Done() chan struct{}
}

// verify stream implements Stream interface
var _ Stream = (*stream)(nil)
var _ Writer = (*writer)(nil)
var _ Reader = (*reader)(nil)

type Stream2 struct {
	elementsCh chan Element
	stopCh chan struct{}
}

func (s *Stream2) Writer() *writer {
	return &writer{stream: s}
}

func (s *Stream2) Reader() *reader {
	return &reader{stream: s}
}


type reader struct {
	name string
	stream *Stream2
	completing bool
}

func (s *reader) Read() <-chan Element {
	return s.stream.elementsCh
}

func (s *reader) Completing() bool {
	return s.completing
}



func (s *reader) Complete() {
	close(s.stream.stopCh)
	s.completing = true
}

type writer struct {
	name string
	stream *Stream2
}

func (w *writer) Done() chan struct{} {
	return w.stream.stopCh
}

func (w *writer) Send(element Element)  {

	select {
	case <- w.stream.stopCh:
		return
	default:
	}

	select {
	case <- w.stream.stopCh:
		return
	case w.stream.elementsCh <- element:
	}
}
func NewStream() * Stream2 {
	return &Stream2{
		elementsCh: make(chan Element),
		stopCh:     make(chan struct{}),
	}
}

