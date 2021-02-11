package sink

import (
	"fmt"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/stream"
)

// verify headOptionStageLogic implements stream.SinkStageLogic interface
var _ stream.SinkStageLogic = (*headOptionStageLogic)(nil)

var HeadOfEmptyStream = fmt.Errorf("head of empty stream")

type headOptionStageLogic struct {
	promise *core.Promise
}

func (h *headOptionStageLogic) OnUpstreamStart(_ stream.SinkStageActions) {

}

func (h *headOptionStageLogic) OnUpstreamReceive(element stream.Element, actions stream.SinkStageActions) {
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

func (h *headOptionStageLogic) OnUpstreamFinish(_ stream.SinkStageActions) {
	if !h.promise.IsCompleted() {
		h.promise.TrySuccess(core.None())
	}
}

func headOptionStage() stream.SinkStage {
	return stream.Sink(createHeadOptionLogic)
}

func createHeadOptionLogic(attributes *stream.StageAttributes) (stream.SinkStageLogic, *core.Promise) {
	promise := core.NewPromise()
	return &headOptionStageLogic{
		promise: promise,
	}, promise
}
