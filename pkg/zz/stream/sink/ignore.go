package sink

import (
	"context"
	"fmt"

	"github.com/nicholasham/piper/pkg/streamold"
)

var Ignored error = fmt.Errorf("output ignored")

// verify ignoredCollector implements CollectorLogic interface
var _ CollectorLogic = (*ignoredCollector)(nil)

type ignoredCollector struct {
}

func (h *ignoredCollector) Start(ctx context.Context, actions CollectActions) {
}

func (h *ignoredCollector) Collect(ctx context.Context, element streamold.Element, actions CollectActions) {
	actions.FailStage(Ignored)
}

func (h *ignoredCollector) End(ctx context.Context, actions CollectActions) {
}

func ignore() CollectorLogic {
	return &headCollector{}
}
