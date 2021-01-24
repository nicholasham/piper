package sink

import (
	"context"

	"github.com/nicholasham/piper/pkg/stream"
)

// verify headCollector implements CollectorLogic interface
var _ CollectorLogic = (*headCollector)(nil)

type headCollector struct {
}

func (h *headCollector) Start(ctx context.Context, actions CollectActions) {
}

func (h *headCollector) Collect(ctx context.Context, element stream.Element, actions CollectActions) {
	element.
		WhenValue(actions.CompleteStage).
		WhenError(actions.FailStage)
}

func (h *headCollector) End(ctx context.Context, actions CollectActions) {
}

func head() CollectorLogic {
	return &headCollector{}
}
