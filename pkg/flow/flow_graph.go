package flow

type FlowGraph struct {
	stage FlowStage
}

func (g *FlowGraph) WithAttributes(attributes ...Attribute) *FlowGraph {
	return FromFlow(g.stage.WithAttributes(attributes...))
}

func FromFlow(stage FlowStage) *FlowGraph {
	return &FlowGraph{
		stage: stage,
	}
}
