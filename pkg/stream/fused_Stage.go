package stream

import "context"

// verify fusedStage implements FlowStage interface
var _ FlowStage = (*fusedStage)(nil)

type fusedStage struct {
	fromStage SourceStage
	toStage   FlowStage
}

func (c *fusedStage) WireTo(stage OutputStage) {
	c.WireTo(stage)
}

func (c *fusedStage) With(options ...StageOption) Stage {
	return NewFusedFlow(c.fromStage, c.toStage.With(options...).(FlowStage))
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

func NewFusedFlow(fromStage SourceStage, toStage FlowStage) FlowStage {
	toStage.WireTo(fromStage)
	return &fusedStage{
		fromStage: fromStage,
		toStage:   toStage,
	}
}
