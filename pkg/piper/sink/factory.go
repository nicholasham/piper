package sink

import (
	"github.com/nicholasham/piper/pkg/piper"
)

func Head(attributes ...piper.StageAttribute) *piper.SinkGraph {
	return piper.SinkFrom(CollectorSink(head(), append(attributes, piper.Name("HeadSink"))))
}

func Ignore(attributes ...piper.StageAttribute) *piper.SinkGraph {
	return piper.SinkFrom(CollectorSink(ignore(), append(attributes, piper.Name("IgnoreSink"))))
}

func List(attributes ...piper.StageAttribute) *piper.SinkGraph {
	return piper.SinkFrom(CollectorSink(list(), append(attributes, piper.Name("ListSink"))))
}
