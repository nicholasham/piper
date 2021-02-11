package sink

import (
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/stream"
)

// verify headOptionStageLogic implements stream.SinkStageLogic interface
var _ stream.SinkStageLogic = (*ignoreSinkStageLogic)(nil)

type ignoreSinkStageLogic struct {
	promise *core.Promise
}

func (i *ignoreSinkStageLogic) OnUpstreamStart(_ stream.SinkStageActions) {
}

func (i *ignoreSinkStageLogic) OnUpstreamReceive(element stream.Element, actions stream.SinkStageActions) {
	element.WhenError(func(err error) {
		actions.FailStage(err)
		i.promise.TryFailure(err)
	})
}

func (i *ignoreSinkStageLogic) OnUpstreamFinish(actions stream.SinkStageActions) {
	i.promise.TrySuccess(stream.NotUsed)
}

func ignoreStage() stream.SinkStage {
	return stream.Sink(createIgnoreStageLogic)
}

func createIgnoreStageLogic(attributes *stream.StageAttributes) (stream.SinkStageLogic, *core.Promise) {
	promise := core.NewPromise()
	return &ignoreSinkStageLogic{
		promise: promise,
	}, promise
}