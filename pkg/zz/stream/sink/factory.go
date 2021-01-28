package sink

import (
	"github.com/nicholasham/piper/pkg/zz/stream"
)

func Head(attributes ...stream.StageOption) *stream.SinkGraph {
	return Collector("HeadSink", head(), attributes...)
}

func Ignore(attributes ...stream.StageOption) *stream.SinkGraph {
	return Collector("IgnoreSink", ignore(), attributes...)
}

func List(attributes ...stream.StageOption) *stream.SinkGraph {
	return Collector("ListSink", list(), attributes...)
}

func Collector(name string, logic CollectorLogic, attributes ...stream.StageOption) *stream.SinkGraph {
	return stream.SinkFrom(CollectorSink(name, logic, attributes))
}
