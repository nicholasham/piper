package stream

import (
	"context"
	"fmt"
	"github.com/gammazero/workerpool"
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
	CompleteStage()

	StageIsCompleted() bool
}

type FlowStageLogicFactory func(attributes *StageAttributes) FlowStageLogic

// verify flowStage implements UpstreamStage interface
var _ UpstreamStage = (*flowStage)(nil)

type flowStage struct {
	attributes    *StageAttributes
	upstreamStage UpstreamStage
	factory       FlowStageLogicFactory
}

func (s *flowStage) Named(name string) Stage {
	return s.With(Name(name))
}

func (s *flowStage) With(options ...StageOption) Stage {
	return &flowStage{
		attributes:    s.attributes.With(options...),
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
	reader, inputFuture := s.upstreamStage.Open(ctx, KeepRight)
	go func() {
		fmt.Println("running " + s.attributes.Name)
		writer := outputStream.Writer()
		logic := s.factory(s.attributes)
		wp := s.createWorkerPool(logic)
		actions := s.newActions(reader, writer)
		logic.OnUpstreamStart(actions)

		for element := range reader.Read() {

			select {
			case <-writer.Done():
				fmt.Println(fmt.Sprintf("Stage done %v", s.attributes.Name))
				reader.Complete()
				break
			default:

			}

			select {
			case <-ctx.Done():
				outputPromise.TryFailure(ctx.Err())
				reader.Complete()
				break
			case <-writer.Done():
				fmt.Println(fmt.Sprintf("Stage done %v", s.attributes.Name))
				reader.Complete()
				break
			default:
			}

			if !reader.Completing() {
				submitToPoolInClosure := func(element Element, actions FlowStageActions) func() {
					return func() {
						logic.OnUpstreamReceive(element, actions)
					}
				}
				wp.Submit(submitToPoolInClosure(element, actions))
			}
		}
		fmt.Println("stopping " + s.attributes.Name)
		wp.StopWait()
		logic.OnUpstreamFinish(actions)
		if !outputPromise.IsCompleted() {
			outputPromise.TrySuccess(NotUsed)
		}

		fmt.Println("stopped " + s.attributes.Name)
	}()

	return outputStream.Reader(), mat(inputFuture, outputPromise.Future())

}

func (s *flowStage) createWorkerPool(logic FlowStageLogic) *workerpool.WorkerPool {
	if logic.SupportsParallelism() {
		return workerpool.New(1)
	}
	return workerpool.New(s.attributes.Parallelism)
}

func (s *flowStage) newActions(inputStream Reader, outputStream Writer) FlowStageActions {
	return &flowStageActions{reader: inputStream, writer: outputStream}
}

// verify flowStageActions implements FlowStageActions interface
var _ FlowStageActions = (*flowStageActions)(nil)

type flowStageActions struct {
	logger Logger
	reader Reader
	writer Writer
}

func (f *flowStageActions) StageIsCompleted() bool {
	return f.reader.Completing()
}

func (f *flowStageActions) SendError(cause error) {
	f.writer.Send(Error(cause))
}

func (f *flowStageActions) SendValue(value interface{}) {
	f.writer.Send(Value(value))
}

func (f *flowStageActions) FailStage(cause error) {
	f.logger.Error(cause, "failed stage because")
	f.reader.Complete()
}

func (f *flowStageActions) CompleteStage() {
	f.reader.Complete()
}

func Flow(factory FlowStageLogicFactory) FlowStage {
	return &flowStage{
		attributes:    DefaultStageAttributes,
		upstreamStage: nil,
		factory:       factory,
	}
}
