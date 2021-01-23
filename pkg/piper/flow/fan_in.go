package flow

import (
	"context"
	"github.com/nicholasham/piper/pkg/piper"
)

// verify fanInFlowStage implements piper.FlowStage interface
var _ piper.FlowStage = (*fanInFlowStage)(nil)

type FanInStrategyFactory func(option ...piper.Option) FanInStrategy

type FanInStrategy func(ctx context.Context, inlets []*piper.Inlet, outlet *piper.Outlet)

type fanInFlowStage struct {
	inlets []*piper.Inlet
	outlet *piper.Outlet
	name   string
	fanIn  FanInStrategy
}

func (receiver *fanInFlowStage) Name() string {
	return receiver.name
}

func (receiver *fanInFlowStage) Run(ctx context.Context) {
	receiver.fanIn(ctx, receiver.inlets, receiver.outlet)
}

func (receiver *fanInFlowStage) Outlet() *piper.Outlet {
	return receiver.outlet
}

func (receiver *fanInFlowStage) Wire(stage piper.SourceStage) {
	inlet := piper.NewInletOld(stage.Name()).WireTo(stage.Outlet())
	receiver.inlets = append(receiver.inlets, inlet)
}

func fanInFlow(name string, stages []piper.SourceStage, strategy FanInStrategy, options ...piper.Option) *fanInFlowStage {

	flow := fanInFlowStage{
		outlet: piper.NewOutletOld(name, options...),
		name:   name,
		fanIn:  strategy,
	}

	for _, stage := range stages {
		flow.Wire(stage)
	}

	return &flow
}
