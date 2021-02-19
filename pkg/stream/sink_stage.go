package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
	"sync"
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
	factory       SinkStageLogicFactory
}

func (s *sinkStage) Named(name string) Stage {
	return s.With(Name(name))
}

func (s *sinkStage) With(options ...StageOption) Stage {
	return &sinkStage{
		attributes:    s.attributes.With(options...),
		upstreamStage: s.upstreamStage,
		factory:       s.factory,
	}
}

func (s *sinkStage) WireTo(stage UpstreamStage) SinkStage {
	s.upstreamStage = stage
	return s
}

func (s *sinkStage) Run(ctx context.Context, wg *sync.WaitGroup, combine MaterializeFunc) *core.Future {
	upstreamReceiver, inputFuture := s.upstreamStage.Open(ctx, wg, combine)
	logic, outputPromise := s.factory(s.attributes)
	wg.Add(1)
	go func() {
		actions := s.newActions(upstreamReceiver)
		logic.OnUpstreamStart(actions)
		defer func() {
			wg.Done()
		}()
		for element := range upstreamReceiver.Receive() {

			select {
			case <-ctx.Done():
				outputPromise.TryFailure(ctx.Err())
				upstreamReceiver.Done()
				return
			default:
			}

			logic.OnUpstreamReceive(element, actions)

		}
		logic.OnUpstreamFinish(actions)
	}()
	return combine(inputFuture, outputPromise.Future())
}

func (s *sinkStage) newActions(upstreamReceiver *Receiver) SinkStageActions {
	return &sinkStageActions{
		logger:           s.attributes.Logger,
		upstreamReceiver: upstreamReceiver,
	}
}

// verify sinkStageActions implements SinkStageActions interface
var _ SinkStageActions = (*sinkStageActions)(nil)

type sinkStageActions struct {
	logger           Logger
	upstreamReceiver *Receiver
}

func (s *sinkStageActions) FailStage(cause error) {
	s.logger.Error(cause, "failed stage because")
	s.upstreamReceiver.Done()
}

func (s *sinkStageActions) CompleteStage() {
	s.upstreamReceiver.Done()
}

func Sink(factory SinkStageLogicFactory) SinkStage {
	return &sinkStage{
		attributes:    DefaultStageAttributes,
		upstreamStage: nil,
		factory:       factory,
	}
}
