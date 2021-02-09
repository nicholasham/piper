package experiment

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
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
	// Completes the stage
	CompleteStage()
}

type SinkStageLogicFactory func(attributes *StageAttributes) (SinkStageLogic, *core.Promise)


// verify sinkStage implements SinkStage interface
var _ SinkStage = (*sinkStage)(nil)

type sinkStage struct {
	attributes    *StageAttributes
	upstreamStage UpstreamStage
	factory SinkStageLogicFactory
	mapMaterializedValue MapMaterializedValueFunc
}

func (s *sinkStage) With(options ...StageOption) Stage {
	return &sinkStage{
		attributes:    s.attributes.Apply(options...),
		upstreamStage: s.upstreamStage,
		factory:       s.factory,
		mapMaterializedValue: s.mapMaterializedValue,
	}
}

func (s *sinkStage) WireTo(stage UpstreamStage) SinkStage {
	s.upstreamStage= stage
	return s
}

func (s *sinkStage) Run(ctx context.Context, combine MaterializeFunc) *core.Future {
	inputReader, inputFuture := s.upstreamStage.Open(ctx, combine)
	logic, outputPromise := s.factory(s.attributes)
	go func() {
		actions  := s.newActions(inputReader)
		logic.OnUpstreamStart(actions)
		for element := range inputReader.Elements() {

			select {
			case <-ctx.Done():
				outputPromise.TryFailure(ctx.Err())
				inputReader.Complete()
			default:
			}

			if !inputReader.Completing() {
				logic.OnUpstreamReceive(element, actions )
			}
		}
		logic.OnUpstreamFinish(actions)
	}()
	return combine(inputFuture, outputPromise.Future())
}

func (s *sinkStage) newActions(reader StreamReader) SinkStageActions {
	return & sinkStageActions{
		logger:       s.attributes.Logger,
		inputStream:  reader,
	}
}


// verify sinkStageActions implements SinkStageActions interface
var _ SinkStageActions = (*sinkStageActions)(nil)

type sinkStageActions struct {
	mapMaterialisedValue MapMaterializedValueFunc
	logger Logger
	inputStream StreamReader
}

func (s *sinkStageActions) FailStage(cause error) {
	s.logger.Error(cause, "failed stage because")
	s.inputStream.Complete()
}

func (s *sinkStageActions) CompleteStage() {
	s.inputStream.Complete()
}

func Sink(factory SinkStageLogicFactory) SinkStage {
	return &sinkStage{
		attributes:    DefaultStageAttributes,
		upstreamStage: nil,
		factory:       factory,
	}
}


