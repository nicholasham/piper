package sink

import (
	. "github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/experiment"
)

// verify headOptionStageLogic implements experiment.SinkStageLogic interface
var _ experiment.SinkStageLogic = (*headOptionStageLogic)(nil)

type headOptionStageLogic struct {
	head    Optional
	promise *Promise
}

func (h *headOptionStageLogic) OnUpstreamStart(actions experiment.SinkStageActions) {

}

func (h *headOptionStageLogic) OnUpstreamReceive(element experiment.Element, actions experiment.SinkStageActions) {
	element.
		WhenValue(func(value interface{}) {
			h.head = Some(value)
			actions.CompleteStage(value)
		}).
		WhenError(actions.FailStage)
}

func (h *headOptionStageLogic) OnUpstreamFinish(actions experiment.SinkStageActions) {
	h.head.
		IfSome(h.promise.TrySuccess).
		IfNone(func() {
			h.promise.TryFailure()
		})
}

func HeadSink() experiment.SinkStage {
	return experiment.Sink(headFactory())
}

func headFactory() experiment.SinkStageLogicFactory {
	return func(attributes *experiment.StageAttributes) (experiment.SinkStageLogic, *Promise) {
		promise := NewPromise()
		return &headOptionStageLogic{
			promise: promise,
			head:    None(),
		}, promise
	}
}
