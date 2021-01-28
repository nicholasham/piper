package flow

import (
	"context"

	"github.com/gammazero/workerpool"
	"github.com/nicholasham/piper/pkg/streamold"
)

// verify operatorFlowStage implements stream.FlowStage interface
var _ streamold.FlowStage = (*operatorFlowStage)(nil)


type OperatorLogicFactory func(options streamold.StageOption)  OperatorLogic

type OperatorLogic interface {
	SupportsParallelism() bool
	Start(actions OperatorActions)
	Apply(element streamold.Element, actions OperatorActions)
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
type SendElement func(element streamold.Element)
type OnPush func(element streamold.Element, actions OperatorActions)

type operatorFlowStage struct {
	name        string
	logger      streamold.Logger
	parallelism int
	inlet       *streamold.Inlet
	outlet      *streamold.Outlet
	operator    OperatorLogic
	decider     streamold.Decider
}

func (receiver *operatorFlowStage) Name() string {
	return receiver.name
}

func (receiver *operatorFlowStage) Run(ctx context.Context) {
	go func(ctx context.Context, parallelism int, operator OperatorLogic, inlet *streamold.Inlet, outlet *streamold.Outlet) {
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
				outlet.Send(streamold.Error(ctx.Err()))
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
			receiver.outlet.Send(streamold.Error(cause))
		},
		pushValue: func(value interface{}) {
			receiver.outlet.Send(streamold.Value(value))
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

func (receiver *operatorFlowStage) Push(element streamold.Element, actions OperatorActions) func() {
	return func() {
		receiver.operator.Apply(element, actions)
	}
}

func (receiver *operatorFlowStage) Outlet() *streamold.Outlet {
	return receiver.outlet
}

func (receiver *operatorFlowStage) Wire(stage streamold.SourceStage) {
	receiver.inlet.WireTo(stage.Outlet())
}

func OperatorFlow(name string, operator OperatorLogic, options ...streamold.StageOption) streamold.FlowStage {

	stageOptions := streamold.DefaultStageOptions.
		Apply(streamold.Name(name)).
		Apply(options...)

	if !operator.SupportsParallelism() {
		stageOptions.Apply(streamold.Parallelism(1))
	}

	return &operatorFlowStage{
		name:        stageOptions.Name,
		logger:      stageOptions.Logger,
		parallelism: stageOptions.Parallelism,
		decider:     stageOptions.Decider,
		operator:    operator,
		inlet:       streamold.NewInlet(stageOptions),
		outlet:      streamold.NewOutlet(stageOptions),
	}
}
