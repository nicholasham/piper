package sink

import (
	"fmt"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/experiment"
)

// verify headOptionStageLogic implements experiment.SinkStageLogic interface
var _ experiment.SinkStageLogic = (*headOptionStageLogic)(nil)

var HeadOfEmptyStream = fmt.Errorf("head of empty stream")

type headOptionStageLogic struct {
	promise *core.Promise
}

func (h *headOptionStageLogic) OnUpstreamStart(_ experiment.SinkStageActions) {

}

func (h *headOptionStageLogic) OnUpstreamReceive(element experiment.Element, actions experiment.SinkStageActions) {
	element.
		WhenValue(func(value interface{}) {
			h.promise.TrySuccess(core.Some(value))
			actions.CompleteStage()
		}).
		WhenError(func(err error) {
			h.promise.TryFailure(err)
			actions.FailStage(err)
	})
}

func (h *headOptionStageLogic) OnUpstreamFinish(_ experiment.SinkStageActions) {
	if !h.promise.IsCompleted() {
		h.promise.TrySuccess(core.None())
	}
}

func HeadOptionStage() experiment.SinkStage {
	return experiment.Sink(createHeadOptionLogic)
}

func createHeadOptionLogic(attributes *experiment.StageAttributes) (experiment.SinkStageLogic, *core.Promise) {
	promise := core.NewPromise()
	return &headOptionStageLogic{
		promise: promise,
	}, promise
}
