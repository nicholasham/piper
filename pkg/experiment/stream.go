package experiment


type Stream interface {
	Reader() StreamReader
	Writer() StreamWriter
}


type StreamReader interface {
	Elements() <- chan Element
	Complete()
}

type StreamWriter interface {
	Close()
	SendValue(value interface{})
	SendError(value interface{})
}

// verify stream implements Stream interface
var _ Stream = (*stream)(nil)


type stream struct {
}

func (s *stream) Reader() StreamReader {
	panic("implement me")
}

func (s *stream) Writer() StreamWriter {
	panic("implement me")
}

func  NewStream() Stream {
	return &stream{

	}
}

