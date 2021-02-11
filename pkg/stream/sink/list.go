package sink

import (
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/stream"
)

// verify sliceStageLogic implements stream.SinkStageLogic interface
var _ stream.SinkStageLogic = (*sliceStageLogic)(nil)


type sliceStageLogic struct {
	values[] interface{}
	promise *core.Promise
}

func (s *sliceStageLogic) OnUpstreamStart(actions stream.SinkStageActions) {
}

func (s *sliceStageLogic) OnUpstreamReceive(element stream.Element, actions stream.SinkStageActions) {
	element.WhenValue(func(value interface{}) {
		s.values = append(s.values, value)
	}).WhenError(func(err error) {
		s.promise.TryFailure(err)
		actions.FailStage(err)
	})
}

func (s *sliceStageLogic) OnUpstreamFinish(actions stream.SinkStageActions) {
	if !s.promise.IsCompleted() {
		s.promise.TrySuccess(s.values)
	}
}


func sliceStage() stream.SinkStage {
	return stream.Sink(createSliceStageLogic)
}

func createSliceStageLogic(attributes *stream.StageAttributes) (stream.SinkStageLogic, *core.Promise) {
	promise := core.NewPromise()
	return &sliceStageLogic{
		promise: promise,
		values: []interface{}{},
	}, promise
}