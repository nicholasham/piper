package sink

import (
	"github.com/nicholasham/piper/pkg/experiment"
)

func Head() *experiment.SinkGraph{
	return experiment.FromSink(HeadSink())
}
