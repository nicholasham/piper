package stream

import (
	"context"
)

// verify alsoToFlowStage implements FlowStage interface
var _ FlowStage = (*alsoToFlowStage)(nil)

type alsoToFlowStage struct {
	name      string
	inlet           *Inlet
	defaultOutlet   *Outlet
	diversionOutlet *Outlet
	f               PredicateFunc
}

func (receiver *alsoToFlowStage) Name() string {
	return receiver.name
}

func (receiver *alsoToFlowStage) Run(ctx context.Context) {
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
				outlet.Send(element)
				diversion.Send(element)
			}
		}
	}(receiver.inlet, receiver.f, receiver.diversionOutlet, receiver.defaultOutlet)
}

func (receiver *alsoToFlowStage) Outlet() *Outlet {
	return receiver.defaultOutlet
}

func (receiver *alsoToFlowStage) Wire(stage SourceStage) {
	receiver.inlet.WireTo(stage.Outlet())
}

func alsoTo(source SourceStage, sink SinkStage, options []StageOption) FlowStage {
	state := NewStageState("AlsoToFlow", options...)
	flow := &alsoToFlowStage{
		name:      state.Name,
		inlet:           NewInlet(state),
		defaultOutlet:   NewOutlet(state),
		diversionOutlet: NewOutlet(NewStageState(state.Name + "-Also")),
	}
	flow.Wire(source)
	sink.Inlet().WireTo(flow.diversionOutlet)
	return flow
}
