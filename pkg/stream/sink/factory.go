package sink

import (
	"github.com/nicholasham/piper/pkg/stream"
)

func Head(attributes ...stream.StageAttribute) *stream.SinkGraph {
	return stream.SinkFrom(CollectorSink(head(), append(attributes, stream.Name("HeadSink"))))
}

func Ignore(attributes ...stream.StageAttribute) *stream.SinkGraph {
	return stream.SinkFrom(CollectorSink(ignore(), append(attributes, stream.Name("IgnoreSink"))))
}

func List(attributes ...stream.StageAttribute) *stream.SinkGraph {
	return stream.SinkFrom(CollectorSink(list(), append(attributes, stream.Name("ListSink"))))
}
