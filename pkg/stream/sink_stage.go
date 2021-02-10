package stream

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

// verify SinkStage implements stream.SinkStage interface
var _ SinkStage = (*sinkStage)(nil)
var _ SinkStageActions = (*sinkStage)(nil)

type sinkStage struct {
	attributes *StageAttributes
	inlet      *Inlet
	promise    *core.Promise
	factory    SinkStageLogicFactory
}

func (s *sinkStage) FailStage(cause error) {
	s.attributes.Logger.Error(cause, "failed stage because")
	s.promise.TryFailure(cause)
	s.inlet.Complete()
}

func (s *sinkStage) CompleteStage(value interface{}) {
	s.inlet.Complete()
	s.promise.TrySuccess(value)
}

func (s *sinkStage) Name() string {
	return s.attributes.Name
}

func (s *sinkStage) Run(ctx context.Context) {
	go func() {
		logic := s.factory(s.attributes)
		logic.OnUpstreamStart(s)
		for element := range s.inlet.In() {

			select {
			case <-ctx.Done():
				s.promise.TryFailure(ctx.Err())
				s.inlet.Complete()
			default:
			}

			if !s.inlet.CompletionSignaled() {
				logic.OnUpstreamReceive(element, s)
			}

		}
		logic.OnUpstreamFinish(s)
	}()
}

func (s *sinkStage) With(options ...StageOption) Stage {
	attributes := s.attributes.Apply(options...)
	return &sinkStage{
		attributes: attributes,
		inlet:      NewInlet(attributes),
		factory:    s.factory,
	}
}

func (s *sinkStage) WireTo(stage OutputStage) SinkStage {
	s.inlet.WireTo(stage.Outlet())
	return s
}

func (s *sinkStage) Result() Future {
	return s.promise.Future()
}

func Sink(factory SinkStageLogicFactory) SinkStage {
	attributes := DefaultStageAttributes.Apply(Name("SinkStage"))
	return &sinkStage{
		attributes: attributes,
		factory:    factory,
		inlet:      NewInlet(attributes),
		promise:    core.NewPromise(),
	}
}
