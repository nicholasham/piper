package sink

import (
	"context"

	"github.com/nicholasham/piper/pkg/streamold"
)

// verify listCollector implements CollectorLogic interface
var _ CollectorLogic = (*listCollector)(nil)

type listCollector struct {
	buffer streamold.ConcurrentBuffer
}

func (l *listCollector) Start(ctx context.Context, actions CollectActions) {

}

func (l *listCollector) Collect(ctx context.Context, element streamold.Element, actions CollectActions) {
	element.
		WhenValue(l.buffer.Append).
		WhenError(actions.FailStage)
}

func (l *listCollector) End(ctx context.Context, actions CollectActions) {
	actions.CompleteStage(l.buffer.Result())
}

func list() CollectorLogic {
	return &listCollector{
		buffer: streamold.NewConcurrentBuffer(),
	}
}
