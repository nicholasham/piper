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
	head core.Optional
}

func (h *headOptionStageLogic) OnUpstreamStart(_ experiment.SinkStageActions) {

}

func (h *headOptionStageLogic) OnUpstreamReceive(element experiment.Element, actions experiment.SinkStageActions) {
	element.
		WhenValue(func(value interface{}) {
			actions.CompleteStage()
		}).
		WhenError(actions.FailStage)
}

func (h *headOptionStageLogic) OnUpstreamFinish(_ experiment.SinkStageActions) {

}

func HeadSink() experiment.SinkStage {
	return experiment.Sink(headLogic)
}

func headLogic(attributes *experiment.StageAttributes) (experiment.SinkStageLogic, *core.Promise) {
	promise := core.NewPromise()
	return &headOptionStageLogic{
		head: core.None(),
	}
}
