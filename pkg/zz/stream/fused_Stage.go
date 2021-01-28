package stream

import "context"

// verify fusedStage implements FlowStage interface
var _ FlowStage = (*fusedStage)(nil)

type fusedStage struct {
	inputStage  SourceStage
	outputStage FlowStage
}

func (c *fusedStage) Name() string {
	return c.outputStage.Name()
}

func (c *fusedStage) Run(ctx context.Context) {
	c.inputStage.Run(ctx)
	c.outputStage.Run(ctx)
}

func (c *fusedStage) Outlet() *Outlet {
	return c.outputStage.Outlet()
}

func (c *fusedStage) Wire(stage SourceStage) {
	c.outputStage.Wire(stage)
}

func NewFusedFlow(inputStage SourceStage, outputStage FlowStage) FlowStage{
	outputStage.Wire(inputStage)
	return &fusedStage{
		inputStage:  inputStage,
		outputStage: outputStage,
	}
}

