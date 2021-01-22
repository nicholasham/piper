package sink

import (
	"context"

	"github.com/form3.tech/piper/pkg/stream"
)

// verify headCollector implements Collector interface
var _ Collector = (*headCollector)(nil)

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

func head() Collector {
	return &headCollector{}
}
