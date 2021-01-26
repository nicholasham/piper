package stream

import (
	"context"
)

// verify diversionFlowStage implements stream.FlowStage interface
var _ FlowStage = (*diversionFlowStage)(nil)

type PredicateFunc func(element Element) bool

type diversionFlowStage struct {
	name string
	inlet           *Inlet
	defaultOutlet   *Outlet
	diversionOutlet *Outlet
	f               PredicateFunc
}

func (receiver *diversionFlowStage) Name() string {
	return receiver.name
}

func (receiver *diversionFlowStage) Run(ctx context.Context) {
	go func(inlet *Inlet, shouldDivert PredicateFunc, diversion *Outlet, outlet *Outlet) {
		defer func() {
			diversion.Close()
			outlet.Close()
		}()
		for element := range receiver.inlet.In() {
			select {
			case <-ctx.Done():
				inlet.Complete()
			case <-outlet.Done():
				inlet.Complete()
			default:

			}

			if element.IsError() {
				outlet.Send(element)
				inlet.Complete()
			}

			if !inlet.CompletionSignaled() {
				if shouldDivert(element) {
					diversion.Send(element)
				} else {
					outlet.Send(element)
				}
			}
		}
	}(receiver.inlet, receiver.f, receiver.diversionOutlet, receiver.defaultOutlet)
}

func (receiver *diversionFlowStage) Outlet() *Outlet {
	return receiver.defaultOutlet
}

func (receiver *diversionFlowStage) Wire(stage SourceStage) {
	receiver.inlet.WireTo(stage.Outlet())
}

func diversion(source SourceStage, sink SinkStage, predicate PredicateFunc, attributes []StageOption) FlowStage {
	state := NewStageState("DiversionFlow", attributes...)
	flow := &diversionFlowStage{
		name:      state.Name,
		inlet:           NewInlet(state),
		defaultOutlet:   NewOutlet(state),
		diversionOutlet: NewOutlet(NewStageState(state.Name + "-Diversion")),
		f:               predicate,
	}
	flow.Wire(source)
	sink.Inlet().WireTo(flow.diversionOutlet)
	return flow
}
