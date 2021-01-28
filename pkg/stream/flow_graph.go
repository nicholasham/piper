package stream

type FlowGraph struct {
	stage FlowStageWithOptions
}

func (g *FlowGraph) WithOptions(options ...StageOption) *FlowGraph {
	return FromFlow(g.stage.WithOptions(options...))
}

func (g *FlowGraph) Via(that *FlowGraph) *FlowGraph {
	return FromFlow(NewFusedFlow(g.stage, that.stage))
}

func FromFlow(stage FlowStageWithOptions) *FlowGraph {
	return &FlowGraph{
		stage: stage,
	}
}
