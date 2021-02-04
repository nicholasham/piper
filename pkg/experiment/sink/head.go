package sink

import (
	"github.com/nicholasham/piper/pkg/experiment"
)

// verify headSinkStage implements experiment.SinkStageLogic interface
var _ experiment.SinkStageLogic = (*headSinkStageLogic)(nil)

type headSinkStageLogic struct {

}

func (h *headSinkStageLogic) OnUpstreamStart(actions experiment.SinkStageActions) {
}

func (h *headSinkStageLogic) OnUpstreamReceive(element experiment.Element, actions experiment.SinkStageActions) {
	element.
		WhenValue(actions.CompleteStage).
		WhenError(actions.FailStage)
}

func (h *headSinkStageLogic) OnUpstreamFinish(actions experiment.SinkStageActions) {
}

func head() experiment.SinkStageLogic {
	return &headSinkStageLogic{

	}
}

func HeadSink() experiment.SinkStage {
	return experiment.Sink(headFactory())
}

func headFactory() experiment.SinkStageLogicFactory {
	return func(attributes *experiment.StageAttributes) experiment.SinkStageLogic {
		return &headSinkStageLogic{
		}
	}
}




