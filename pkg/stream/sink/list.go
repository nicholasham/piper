package sink

import (
	"context"

	"github.com/nicholasham/piper/pkg/stream"
)

// verify listCollector implements Collector interface
var _ Collector = (*listCollector)(nil)

type listCollector struct {
	buffer stream.ConcurrentBuffer
}

func (l *listCollector) Start(ctx context.Context, actions CollectActions) {

}

func (l *listCollector) Collect(ctx context.Context, element stream.Element, actions CollectActions) {
	element.
		WhenValue(l.buffer.Append).
		WhenError(actions.FailStage)
}

func (l *listCollector) End(ctx context.Context, actions CollectActions) {
	actions.CompleteStage(l.buffer.Result())
}

func list() Collector {
	return &listCollector{
		buffer: stream.NewConcurrentBuffer(),
	}
}
