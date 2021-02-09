package sink

import (
	. "github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/experiment"
)

func HeadOption() *experiment.SinkGraph{
	return experiment.FromSink(HeadOptionStage())
}

func Head() *experiment.SinkGraph{
	return HeadOption().
		MapMaterializedValue(func(value Any) Result {
		return value.(Optional).ToResult(func() error {
			return HeadOfEmptyStream
		})
	})
}
