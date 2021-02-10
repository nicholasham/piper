package sink

import "github.com/nicholasham/piper/pkg/old-stream"

func Head() *old_stream.SinkGraph {
	return old_stream.FromSink(HeadSink())
}
