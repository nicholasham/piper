package stream

type FlowGraph struct {
	stage FlowStage
}

func (g *FlowGraph) With(options ...StageOption) *FlowGraph {
	return FromFlow(g.stage.With(options...).(FlowStage))
}

func (g *FlowGraph) Via(that *FlowGraph) *FlowGraph {
	return FromFlow(NewFusedFlow(g.stage, that.stage))
}

func FromFlow(stage FlowStage) *FlowGraph {
	return &FlowGraph{
		stage: stage,
	}
}
