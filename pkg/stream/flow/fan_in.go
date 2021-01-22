package flow

import (
	"context"
	"github.com/nicholasham/piper/pkg/stream"
)

// verify fanInFlowStage implements stream.FlowStage interface
var _ stream.FlowStage = (*fanInFlowStage)(nil)

type FanInStrategyFactory func(option ...stream.Option) FanInStrategy

type FanInStrategy func(ctx context.Context, inlets []*stream.Inlet, outlet *stream.Outlet)

type fanInFlowStage struct {
	inlets []*stream.Inlet
	outlet *stream.Outlet
	name   string
	fanIn  FanInStrategy
}

func (receiver *fanInFlowStage) Name() string {
	return receiver.name
}

func (receiver *fanInFlowStage) Run(ctx context.Context) {
	receiver.fanIn(ctx, receiver.inlets, receiver.outlet)
}

func (receiver *fanInFlowStage) Outlet() *stream.Outlet {
	return receiver.outlet
}

func (receiver *fanInFlowStage) Wire(stage stream.SourceStage) {
	inlet := stream.NewInletOld(stage.Name()).WireTo(stage.Outlet())
	receiver.inlets = append(receiver.inlets, inlet)
}

func fanInFlow(name string, stages []stream.SourceStage, strategy FanInStrategy, options ...stream.Option) *fanInFlowStage {

	flow := fanInFlowStage{
		outlet: stream.NewOutletOld(name, options...),
		name:   name,
		fanIn:  strategy,
	}

	for _, stage := range stages {
		flow.Wire(stage)
	}

	return &flow
}
