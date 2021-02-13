package stream

// https://go101.org/article/channel-closing.html


type Stream2 struct {
	elementsCh chan Element
	stopCh chan struct{}
}

func (s *Stream2) Writer() *SWriter {
	return &SWriter{stream: s}
}

func (s *Stream2) Reader() *SReader {
	return &SReader{stream: s}
}


type SReader struct {
	name string
	stream *Stream2
}

func (s *SReader) Read() chan Element {
	return s.stream.elementsCh
}

func (s *SReader) Close() {
	close(s.stream.stopCh)
}

type SWriter struct {
	name string
	stream *Stream2
}

func (s *SWriter) Write(element Element)  {

	select {
	case <- s.stream.stopCh:
		return
	default:
	}

	select {
	case <- s.stream.stopCh:
		return
	case s.stream.elementsCh <- element:
	}
}

