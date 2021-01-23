package piper

import (
	"context"
)

// verify diversionFlowStage implements piper.FlowStage interface
var _ FlowStage = (*diversionFlowStage)(nil)

type PredicateFunc func(element Element) bool

type diversionFlowStage struct {
	attributes      *StageAttributes
	inlet           *Inlet
	defaultOutlet   *Outlet
	diversionOutlet *Outlet
	f               PredicateFunc
}

func (receiver *diversionFlowStage) Name() string {
	return receiver.attributes.Name
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

func diversion(source SourceStage, sink SinkStage, predicate PredicateFunc, attributes []StageAttribute) FlowStage {
	stageAttributes := NewAttributes("DiversionFlow", attributes...)
	flow := &diversionFlowStage{
		attributes:      stageAttributes,
		inlet:           NewInlet(stageAttributes),
		defaultOutlet:   NewOutlet(stageAttributes),
		diversionOutlet: NewOutlet(NewAttributes(stageAttributes.Name + "-Diversion")),
		f:               predicate,
	}
	flow.Wire(source)
	sink.Inlet().WireTo(flow.diversionOutlet)
	return flow
}
