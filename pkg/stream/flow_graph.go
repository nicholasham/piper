package stream

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
	return FromFlow(NewFusedFlow(g.stage, that.stage))
}

func FromFlow(stage FlowStage) *FlowGraph {
	return &FlowGraph{
		stage: stage,
	}
}

func (g *FlowGraph) Concat(that *FlowGraph) *FlowGraph {
	return ConcatFlows(g, that)
}

func (g *FlowGraph) Interleave(segmentSize int, that *FlowGraph) *FlowGraph {
	return InterleaveFlows(segmentSize, g, that)
}

func (g *FlowGraph) Merge(that *FlowGraph) *FlowGraph {
	return MergeFlows(g, that)
}
