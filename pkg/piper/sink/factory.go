package sink

import (
	"github.com/nicholasham/piper/pkg/piper"
	"github.com/nicholasham/piper/pkg/piper/attribute"
)

func Head(attributes ...attribute.StageAttribute) *piper.SinkGraph {
	return piper.SinkFrom(CollectorSink(head(), append(attributes, attribute.Name("HeadSink"))))
}

func Ignore(attributes ...attribute.StageAttribute) *piper.SinkGraph {
	return piper.SinkFrom(CollectorSink(ignore(), append(attributes, attribute.Name("IgnoreSink"))))
}

func List(attributes ...attribute.StageAttribute) *piper.SinkGraph {
	return piper.SinkFrom(CollectorSink(list(), append(attributes, attribute.Name("ListSink"))))
}
