package experiment

import (
	"github.com/nicholasham/piper/pkg/core"
	"golang.org/x/net/context"
)

type SinkStageLogic interface {
	// Called when starting to receive elements from upstream
	OnUpstreamStart(actions SinkStageActions)
	// Called when an element is received from upstream
	OnUpstreamReceive(element Element, actions SinkStageActions)
	// 	Called when finishing receiving elements from upstream
	OnUpstreamFinish(actions SinkStageActions)
}

type SinkStageActions interface {
	// Fails a stage and logs the cause of failure.
	FailStage(cause error)
	// Completes the stage with a final value
	CompleteStage(value interface{})
}

type SinkStageLogicFactory func(attributes *StageAttributes) SinkStageLogic


// verify headSourceStage implements UpstreamStage interface
//var _ UpstreamStage = (*headSourceStage)(nil)

type sinkStage struct {
	upstreamStage UpstreamStage
	factory SinkStageLogicFactory
}

func (s *sinkStage) Run(ctx context.Context, mat MaterializeFunc) *core.Promise {
	s.factory :=
	inputReader, inputPromise := s.upstreamStage.Open(ctx, mat)
	go func() {
		for element := range inputReader.Elements() {

		}
	}()

	return mat(inputPromise, inputPromise)
}

