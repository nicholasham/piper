package stream

import "context"

// verify compositeStage implements FlowStage interface
var _ FlowStage = (*compositeStage)(nil)

type compositeStage struct {
	inputStage  SourceStage
	outputStage FlowStage
}

func (c *compositeStage) Name() string {
	return c.outputStage.Name()
}

func (c *compositeStage) Run(ctx context.Context) {
	c.inputStage.Run(ctx)
	c.outputStage.Run(ctx)
}

func (c *compositeStage) Outlet() *Outlet {
	return c.outputStage.Outlet()
}

func (c *compositeStage) Wire(stage SourceStage) {
	c.outputStage.Wire(stage)
}

func NewCompositeFlow(inputStage SourceStage, outputStage FlowStage) FlowStage{
	outputStage.Wire(inputStage)
	return &compositeStage{
		inputStage:  inputStage,
		outputStage: outputStage,
	}
}

