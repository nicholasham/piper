package flow

import (
	"context"

	"github.com/gammazero/workerpool"
	"github.com/nicholasham/piper/pkg/stream"
)

// verify operatorFlowStage implements stream.FlowStage interface
var _ stream.FlowStage = (*operatorFlowStage)(nil)


type OperatorLogicFactory func(options stream.StageOption)  OperatorLogic

type OperatorLogic interface {
	SupportsParallelism() bool
	Start(actions OperatorActions)
	Apply(element stream.Element, actions OperatorActions)
	End(actions OperatorActions)
}

type OperatorActions interface {
	PushError(cause error)
	PushValue(value interface{})
	FailStage(cause error)
	CompleteStage()
}

// verify operatorFlowStage implements stream.FlowStage interface
var _ OperatorActions = (*operatorActions)(nil)

type operatorActions struct {
	pushError     func(cause error)
	pushValue     func(value interface{})
	failStage     func(cause error)
	completeStage func()
}

func (o *operatorActions) PushError(cause error) {
	o.pushError(cause)
}

func (o *operatorActions) PushValue(value interface{}) {
	o.pushValue(value)
}

func (o *operatorActions) FailStage(cause error) {
	o.failStage(cause)
}

func (o *operatorActions) CompleteStage() {
	o.completeStage()
}

type CompleteStage func()
type FailStage func(cause error)
type SendElement func(element stream.Element)
type OnPush func(element stream.Element, actions OperatorActions)

type operatorFlowStage struct {
	name        string
	logger      stream.Logger
	parallelism int
	inlet       *stream.Inlet
	outlet      *stream.Outlet
	operator    OperatorLogic
	decider     stream.Decider
}

func (receiver *operatorFlowStage) Name() string {
	return receiver.name
}

func (receiver *operatorFlowStage) Run(ctx context.Context) {
	go func(ctx context.Context, parallelism int, operator OperatorLogic, inlet *stream.Inlet, outlet *stream.Outlet) {
		wp := workerpool.New(parallelism)
		defer func() {
			outlet.Close()
		}()
		actions := receiver.newOperatorActions()
		operator.Start(actions)
		for element := range inlet.In() {

			if !inlet.CompletionSignaled() {
				wp.Submit(receiver.Push(element, actions))
			}

			select {
			case <-ctx.Done():
				outlet.Send(stream.Error(ctx.Err()))
				inlet.Complete()
			case <-outlet.Done():
				inlet.Complete()
			default:
			}

		}
		wp.StopWait()
		operator.End(actions)
	}(ctx, receiver.parallelism, receiver.operator, receiver.inlet, receiver.outlet)
}

func (receiver *operatorFlowStage) newOperatorActions() OperatorActions {
	return &operatorActions{
		pushError: func(cause error) {
			receiver.outlet.Send(stream.Error(cause))
		},
		pushValue: func(value interface{}) {
			receiver.outlet.Send(stream.Value(value))
		},
		failStage: func(cause error) {
			receiver.logger.Error(cause, "failed stage because")
			receiver.inlet.Complete()
		},
		completeStage: func() {
			receiver.inlet.Complete()
		},
	}
}

func (receiver *operatorFlowStage) Push(element stream.Element, actions OperatorActions) func() {
	return func() {
		receiver.operator.Apply(element, actions)
	}
}

func (receiver *operatorFlowStage) Outlet() *stream.Outlet {
	return receiver.outlet
}

func (receiver *operatorFlowStage) Wire(stage stream.SourceStage) {
	receiver.inlet.WireTo(stage.Outlet())
}

func OperatorFlow(name string, operator OperatorLogic, options ...stream.StageOption) stream.FlowStage {

	stageOptions := stream.DefaultStageOptions.
		Apply(stream.Name(name)).
		Apply(options...)

	if !operator.SupportsParallelism() {
		stageOptions.Apply(stream.Parallelism(1))
	}

	return &operatorFlowStage{
		name: stageOptions.Name,
		logger: stageOptions.Logger,
		parallelism: stageOptions.Parallelism,
		decider: stageOptions.Decider,
		operator:   operator,
		inlet:      stream.NewInlet(stageOptions),
		outlet:     stream.NewOutlet(stageOptions),
	}
}
