package sink

import (
	. "github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/stream"
)

func HeadOption() *stream.SinkGraph {
	return stream.FromSink(HeadOptionStage())
}

func Head() *stream.SinkGraph {
	return HeadOption().
		MapMaterializedValue(func(value Any) Result {
			return value.(Optional).ToResult(func() error {
				return HeadOfEmptyStream
			})
		})
}
