package experiment

type FlowGraph struct {
	stage FlowStage
}

func (g *FlowGraph) With(options ...StageOption) *FlowGraph {
	return FromFlow(g.stage.With(options...).(FlowStage))
}

func (g *FlowGraph) Named(name string) *FlowGraph {
	return g.With(Name(name))
}

func (g *FlowGraph) Via(that *FlowGraph) *FlowGraph {
	that.stage.WireTo(g.stage)
	return that
}

func FromFlow(stage FlowStage) *FlowGraph {
	return &FlowGraph{
		stage: stage,
	}
}

