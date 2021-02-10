package sink

import "github.com/nicholasham/piper/pkg/old-stream"

// verify headSinkStage implements old-stream.SinkStageLogic interface
var _ old_stream.SinkStageLogic = (*headSinkStageLogic)(nil)

type headSinkStageLogic struct {
}

func (h *headSinkStageLogic) OnUpstreamStart(actions old_stream.SinkStageActions) {
}

func (h *headSinkStageLogic) OnUpstreamReceive(element old_stream.Element, actions old_stream.SinkStageActions) {
	element.
		WhenValue(actions.CompleteStage).
		WhenError(actions.FailStage)
}

func (h *headSinkStageLogic) OnUpstreamFinish(actions old_stream.SinkStageActions) {
}

func head() old_stream.SinkStageLogic {
	return &headSinkStageLogic{}
}

func HeadSink() old_stream.SinkStage {
	return old_stream.Sink(headFactory())
}

func headFactory() old_stream.SinkStageLogicFactory {
	return func(attributes *old_stream.StageAttributes) old_stream.SinkStageLogic {
		return &headSinkStageLogic{}
	}
}
