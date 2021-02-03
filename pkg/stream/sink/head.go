package sink

import "github.com/nicholasham/piper/pkg/stream"

// verify headSinkStage implements stream.SinkStageLogic interface
var _ stream.SinkStageLogic = (*headSinkStageLogic)(nil)

type headSinkStageLogic struct {

}

func (h *headSinkStageLogic) OnUpstreamStart(actions stream.SinkStageActions) {
}

func (h *headSinkStageLogic) OnUpstreamReceive(element stream.Element, actions stream.SinkStageActions) {
	element.
		WhenValue(actions.CompleteStage).
		WhenError(actions.FailStage)
}

func (h *headSinkStageLogic) OnUpstreamFinish(actions stream.SinkStageActions) {
}

func head() stream.SinkStageLogic {
	return &headSinkStageLogic{

	}
}

func HeadSink() stream.SinkStage {
	return stream.Sink(headFactory())
}

func headFactory() stream.SinkStageLogicFactory {
	return func(attributes *stream.StageAttributes) stream.SinkStageLogic {
		return &headSinkStageLogic{
		}
	}
}




