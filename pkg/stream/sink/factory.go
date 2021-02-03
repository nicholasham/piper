package sink

import "github.com/nicholasham/piper/pkg/stream"

func Head() *stream.SinkGraph{
	return stream.FromSink (HeadSink())
}
