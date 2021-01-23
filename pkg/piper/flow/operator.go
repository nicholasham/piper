package flow

import (
	"context"
	"github.com/gammazero/workerpool"
	"github.com/nicholasham/piper/pkg/piper"
)

// verify operatorFlowStage implements piper.FlowStage interface
var _ piper.FlowStage = (*operatorFlowStage)(nil)

type Operator interface {
	SupportsParallelism() bool
	Start(actions OperatorActions)
	Apply(element piper.Element, actions OperatorActions)
	End(actions OperatorActions)
}

type OperatorActions interface {
	SendDownstream(element piper.Element)
	PushError(cause error)
	PushValue(value interface{})
	FailStage(cause error)
	CompleteStage()
}

// verify operatorFlowStage implements piper.FlowStage interface
var _ OperatorActions = (*operatorActions)(nil)

type operatorActions struct {
	pushError     func(cause error)
	pushValue     func(value interface{})
	failStage     func(cause error)
	completeStage func()
}

func (o *operatorActions) SendDownstream(element piper.Element) {
	element.
		WhenValue(o.pushValue).
		WhenError(o.pushError)
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
type SendElement func(element piper.Element)
type OnPush func(element piper.Element, actions OperatorActions)

type operatorFlowStage struct {
	attributes *piper.StageAttributes
	inlet      *piper.Inlet
	outlet     *piper.Outlet
	operator   Operator
}

func (receiver *operatorFlowStage) Name() string {
	return receiver.attributes.Name
}

func (receiver *operatorFlowStage) Run(ctx context.Context) {
	go func(ctx context.Context, parallelism int, operator Operator, inlet *piper.Inlet, outlet *piper.Outlet) {
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
				outlet.Send(piper.Error(ctx.Err()))
				inlet.Complete()
			case <-outlet.Done():
				inlet.Complete()
			default:
			}

		}
		wp.StopWait()
		operator.End(actions)
	}(ctx, receiver.attributes.Parallelism, receiver.operator, receiver.inlet, receiver.outlet)
}

func (receiver *operatorFlowStage) newOperatorActions() OperatorActions {
	return &operatorActions{
		pushError: func(cause error) {
			receiver.outlet.Send(piper.Error(cause))
		},
		pushValue: func(value interface{}) {
			receiver.outlet.Send(piper.Value(value))
		},
		failStage: func(cause error) {
			receiver.attributes.Logger.Error(cause, "failed stage because")
			receiver.inlet.Complete()
		},
		completeStage: func() {
			receiver.inlet.Complete()
		},
	}
}

func (receiver *operatorFlowStage) Push(element piper.Element, actions OperatorActions) func() {
	return func() {
		receiver.operator.Apply(element, actions)
	}
}

func (receiver *operatorFlowStage) Outlet() *piper.Outlet {
	return receiver.outlet
}

func (receiver *operatorFlowStage) Wire(stage piper.SourceStage) {
	receiver.inlet.WireTo(stage.Outlet())
}

func OperatorFlow(operator Operator, attributes ...piper.StageAttribute) piper.FlowStage {
	if !operator.SupportsParallelism() {
		attributes = append(attributes, piper.Parallelism(1))
	}

	stageAttributes := piper.NewAttributes("HeadSink", attributes...)

	return &operatorFlowStage{
		attributes: stageAttributes,
		operator:   operator,
		inlet:      piper.NewInlet(stageAttributes),
		outlet:     piper.NewOutlet(stageAttributes),
	}
}
