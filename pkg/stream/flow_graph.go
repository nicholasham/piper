package stream

type FlowGraph struct {
	stage FlowStage
}

func (g *FlowGraph) WithOptions(options ...StageOption) *FlowGraph {
	return FromFlow(g.stage.WithOptions(options...))
}

func FromFlow(stage FlowStage) *FlowGraph {
	return &FlowGraph{
		stage: stage,
	}
}
