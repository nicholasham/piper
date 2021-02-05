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
	// Completes the stage with a final value
	CompleteStage(value interface{})
}

type SinkStageLogicFactory func(attributes *StageAttributes) SinkStageLogic


// verify sinkStage implements SinkStage interface
var _ SinkStage = (*sinkStage)(nil)

type sinkStage struct {
	attributes    *StageAttributes
	upstreamStage UpstreamStage
	factory SinkStageLogicFactory
}

func (s *sinkStage) With(options ...StageOption) Stage {
	panic("implement me")
}

func (s *sinkStage) WireTo(stage UpstreamStage) SinkStage {
	s.upstreamStage= stage
	return s
}

func (s *sinkStage) Run(ctx context.Context, mat MaterializeFunc) *core.Promise {
	inputReader, inputPromise := s.upstreamStage.Open(ctx, mat)
	outputPromise := core.NewPromise()
	go func() {
		logic := s.factory(s.attributes)
		actions  := s.newActions(inputReader, outputPromise)
		logic.OnUpstreamStart(actions)
		for element := range inputReader.Elements() {

			select {
			case <-ctx.Done():
				outputPromise.Reject(ctx.Err())
				inputReader.Complete()
			default:
			}

			if !inputReader.Completing() {
				logic.OnUpstreamReceive(element, actions )
			}
		}
		logic.OnUpstreamFinish(actions)
	}()
	return mat(inputPromise, outputPromise)
}

func (s *sinkStage) newActions(reader StreamReader, promise * core.Promise) SinkStageActions {
	return & sinkStageActions{
		logger:       s.attributes.Logger,
		inputStream:  reader,
		promise: promise,
	}
}


// verify sinkStageActions implements SinkStageActions interface
var _ SinkStageActions = (*sinkStageActions)(nil)

type sinkStageActions struct {
	logger Logger
	inputStream StreamReader
	promise *core.Promise
}

func (s *sinkStageActions) FailStage(cause error) {
	s.logger.Error(cause, "failed stage because")
	s.inputStream.Complete()
	s.promise.Reject(cause)
}

func (s *sinkStageActions) CompleteStage(value interface{}) {
	s.inputStream.Complete()
	s.promise.Resolve(value)
}

func Sink(factory SinkStageLogicFactory) SinkStage {
	return &sinkStage{
		attributes:    DefaultStageAttributes,
		upstreamStage: nil,
		factory:       factory,
	}
}


