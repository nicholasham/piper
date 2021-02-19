package stream

import (
	"context"
	"fmt"
	"github.com/gammazero/workerpool"
	"github.com/nicholasham/piper/pkg/core"
	"sync"
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

func (s *flowStage) Open(ctx context.Context, wg *sync.WaitGroup, mat MaterializeFunc) (*Receiver, *core.Future) {
	outputStream := NewStream(s.attributes.Name)
	outputPromise := core.NewPromise()
	receiver, inputFuture := s.upstreamStage.Open(ctx, wg, KeepRight)
	wg.Add(1)
	go func() {
		sender := outputStream.Sender()
		logic := s.factory(s.attributes)
		wp := s.createWorkerPool(logic)
		actions := s.newActions(receiver, sender)
		logic.OnUpstreamStart(actions)

		defer func() {
			wp.StopWait()
			sender.Close()
			wg.Done()
		}()

		for element := range receiver.Receive() {

			select {
			case <-ctx.Done():
				outputPromise.TryFailure(ctx.Err())
				receiver.Done()
				return
			case <-sender.Done():
				fmt.Println(fmt.Sprintf("Stage done %v", s.attributes.Name))
				receiver.Done()
				return
			default:
			}

			submitToPoolInClosure := func(element Element, actions FlowStageActions) func() {
				return func() {
					logic.OnUpstreamReceive(element, actions)
				}
			}
			wp.Submit(submitToPoolInClosure(element, actions))
		}
		wp.StopWait()
		logic.OnUpstreamFinish(actions)
		if !outputPromise.IsCompleted() {
			outputPromise.TrySuccess(NotUsed)
		}
	}()

	return outputStream.Receiver(), mat(inputFuture, outputPromise.Future())

}

func (s *flowStage) createWorkerPool(logic FlowStageLogic) *workerpool.WorkerPool {
	if logic.SupportsParallelism() {
		return workerpool.New(1)
	}
	return workerpool.New(s.attributes.Parallelism)
}

func (s *flowStage) newActions(receiver *Receiver, sender *Sender) FlowStageActions {
	return &flowStageActions{receiver: receiver, sender: sender}
}

// verify flowStageActions implements FlowStageActions interface
var _ FlowStageActions = (*flowStageActions)(nil)

type flowStageActions struct {
	logger   Logger
	receiver *Receiver
	sender   *Sender
}

func (f *flowStageActions) StageIsCompleted() bool {
	return f.sender.IsDone()
}

func (f *flowStageActions) SendError(cause error) {
	f.sender.TrySend(Error(cause))
}

func (f *flowStageActions) SendValue(value interface{}) {
	f.sender.TrySend(Value(value))
}

func (f *flowStageActions) FailStage(cause error) {
	f.logger.Error(cause, "failed stage because")
	f.receiver.Done()
}

func (f *flowStageActions) CompleteStage() {
	f.receiver.Done()
}

func Flow(factory FlowStageLogicFactory) FlowStage {
	return &flowStage{
		attributes:    DefaultStageAttributes,
		upstreamStage: nil,
		factory:       factory,
	}
}
