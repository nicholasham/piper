package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
)



type FlowStageLogic interface {
	SupportsParallelism() bool
	// Called when starting to receive elements from upstream
	OnUpstreamStart(actions FlowStageActions)
	// Called when an element is received from upstream
	OnUpstreamReceive(element Element, actions FlowStageActions)
	// 	Called when finishing receiving elements from upstream
	OnUpstreamFinish(actions FlowStageActions)
}

type FlowStageActions interface {
	// Sends an error downstream
	SendError(cause error)
	// Sends a value downstream
	SendValue(value interface{})
	// Fails a stage on logs the cause of failure.
	FailStage(cause error)
	// Completes the stage with a materialised value
	CompleteStage(value interface{})
}

type FlowStageLogicFactory func(attributes *StageAttributes) FlowStageLogic

// verify flowStage implements UpstreamStage interface
var _ UpstreamStage = (*flowStage)(nil)

type flowStage struct {
	attributes    *StageAttributes
	upstreamStage UpstreamStage
	factory       FlowStageLogicFactory
}

func (s *flowStage) With(options ...StageOption) Stage {
	return &flowStage{
		attributes:    s.attributes.Apply(options...),
		upstreamStage: s.upstreamStage,
		factory:       s.factory,
	}
}

func (s *flowStage) WireTo(stage UpstreamStage) FlowStage {
	s.upstreamStage = stage
	return s
}

func (s *flowStage) Open(ctx context.Context, mat MaterializeFunc) (Reader, *core.Future) {
	outputStream := NewStream()
	outputPromise := core.NewPromise()
	reader, inputFuture :=  s.upstreamStage.Open(ctx, KeepRight)
	go func() {
		writer:= outputStream.Writer()
		defer writer.Close()
		logic := s.factory(s.attributes)
		actions := s.newActions(reader, writer)
		logic.OnUpstreamStart(actions)
		for element := range reader.Elements(){
			select {
			case <-ctx.Done():
				outputPromise.TryFailure(ctx.Err())
				reader.Complete()
			case <-writer.Done():
				reader.Complete()
			default:
			}

			if !reader.Completing() {
				logic.OnUpstreamReceive(element, actions)
			}

		}
		logic.OnUpstreamFinish(actions)
	}()
	return outputStream.Reader(), mat(inputFuture, outputPromise.Future())
}

func (s *flowStage) newActions(inputStream Reader, outputStream Writer) FlowStageActions {
	return & flowStageActions{inputStream: inputStream, outputStream: outputStream}
}

// verify flowStageActions implements FlowStageActions interface
var _ FlowStageActions = (*flowStageActions)(nil)

type flowStageActions struct {
	logger       Logger
	inputStream  Reader
	outputStream Writer
}

func (f *flowStageActions) SendError(cause error) {
	f.outputStream.SendError(cause)
}

func (f *flowStageActions) SendValue(value interface{}) {
	f.outputStream.SendValue(value)
}

func (f *flowStageActions) FailStage(cause error) {
	f.logger.Error(cause, "failed stage because")
	f.inputStream.Complete()
}

func (f *flowStageActions) CompleteStage(value interface{}) {
	f.inputStream.Complete()
}

func Flow(factory FlowStageLogicFactory) FlowStage {
	return &flowStage{
		attributes:    DefaultStageAttributes,
		upstreamStage: nil,
		factory:       factory,
	}
}






