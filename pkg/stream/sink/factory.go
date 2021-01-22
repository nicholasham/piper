package sink

import (
	"github.com/nicholasham/piper/pkg/stream"
	"github.com/nicholasham/piper/pkg/stream/attribute"
)

func Head(attributes ...attribute.StageAttribute) *stream.SinkGraph {
	return stream.SinkFrom(CollectorSink(head(), append(attributes, attribute.Name("HeadSink"))))
}

func Ignore(attributes ...attribute.StageAttribute) *stream.SinkGraph {
	return stream.SinkFrom(CollectorSink(ignore(), append(attributes, attribute.Name("IgnoreSink"))))
}

func List(attributes ...attribute.StageAttribute) *stream.SinkGraph {
	return stream.SinkFrom(CollectorSink(list(), append(attributes, attribute.Name("ListSink"))))
}
