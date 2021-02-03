package stream

import "context"

// verify compositeStage implements FlowStage interface
var _ FlowStage = (*compositeStage)(nil)

type compositeStage struct {
	fromStage SourceStage
	toStage   FlowStage
}

func (c *compositeStage) WireTo(stage OutputStage) FlowStage {
	c.WireTo(stage)
	return c
}

func (c *compositeStage) With(options ...StageOption) Stage {
	return CompositeFlow(c.fromStage, c.toStage.With(options...).(FlowStage))
}

func (c *compositeStage) Name() string {
	return c.toStage.Name()
}

func (c *compositeStage) Run(ctx context.Context) {
	c.fromStage.Run(ctx)
	c.toStage.Run(ctx)
}

func (c *compositeStage) Outlet() *Outlet {
	return c.toStage.Outlet()
}

func CompositeFlow(fromStage SourceStage, toStage FlowStage) FlowStage {
	toStage.WireTo(fromStage)
	return &compositeStage{
		fromStage: fromStage,
		toStage:   toStage,
	}
}
