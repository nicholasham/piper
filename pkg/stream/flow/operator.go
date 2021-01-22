package flow

import (
	"context"
	"github.com/nicholasham/piper/pkg/stream"
	"github.com/nicholasham/piper/pkg/stream/attribute"

	"github.com/gammazero/workerpool"

)

// verify operatorFlowStage implements stream.FlowStage interface
var _ stream.FlowStage = (*operatorFlowStage)(nil)

type Operator interface {
	SupportsParallelism() bool
	Start(actions OperatorActions)
	Apply(element stream.Element, actions OperatorActions)
	End(actions OperatorActions)
}

type OperatorActions interface {
	SendDownstream(element stream.Element)
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

func (o *operatorActions) SendDownstream(element stream.Element) {
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
type SendElement func(element stream.Element)
type OnPush func(element stream.Element, actions OperatorActions)

type operatorFlowStage struct {
	attributes *attribute.StageAttributes
	inlet      *stream.Inlet
	outlet     *stream.Outlet
	operator   Operator
}

func (receiver *operatorFlowStage) Name() string {
	return receiver.attributes.Name
}

func (receiver *operatorFlowStage) Run(ctx context.Context) {
	go func(ctx context.Context, parallelism int, operator Operator, inlet *stream.Inlet, outlet *stream.Outlet) {
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
	}(ctx, receiver.attributes.Parallelism, receiver.operator, receiver.inlet, receiver.outlet)
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
			receiver.attributes.Logger.Error(cause, "failed stage because")
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

func OperatorFlow(operator Operator, attributes ...attribute.StageAttribute) stream.FlowStage {
	if !operator.SupportsParallelism() {
		attributes = append(attributes, attribute.Parallelism(1))
	}

	stageAttributes := attribute.Default("HeadSink", attributes...)

	return &operatorFlowStage{
		attributes: stageAttributes,
		operator:   operator,
		inlet:      stream.NewInlet(stageAttributes),
		outlet:     stream.NewOutlet(stageAttributes),
	}
}
