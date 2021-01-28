package stream

import "context"

// verify fusedStage implements FlowStage interface
var _ FlowStage = (*fusedStage)(nil)

type fusedStage struct {
	fromStage SourceStage
	toStage   FlowStageWithOptions
}

func (c *fusedStage) Inlet() *Inlet {
	return c.toStage.Inlet()
}

func (c *fusedStage) WithOptions(options ...StageOption) FlowStageWithOptions {
	return NewFusedFlow(c.fromStage, c.toStage.WithOptions(options...))
}

func (c *fusedStage) Name() string {
	return c.toStage.Name()
}

func (c *fusedStage) Run(ctx context.Context) {
	c.fromStage.Run(ctx)
	c.toStage.Run(ctx)
}

func (c *fusedStage) Outlet() *Outlet {
	return c.toStage.Outlet()
}


func NewFusedFlow(fromStage SourceStage, toStage FlowStageWithOptions) FlowStageWithOptions{
	toStage.Inlet().WireTo(fromStage.Outlet())
	return &fusedStage{
		fromStage: fromStage,
		toStage:   toStage,
	}
}

