package stream

type SourceGraph struct {
	stage SourceStageWithOptions
}

func (g *SourceGraph) WithOptions(options ...StageOption) *SourceGraph {
	return FromSource(g.stage.WithOptions(options...))
}

func (g *SourceGraph) Via(that *FlowGraph) *FlowGraph {
	return FromFlow(NewFusedFlow(g.stage, that.stage))
}

func FromSource(stage SourceStageWithOptions) *SourceGraph {
	return &SourceGraph{
		stage: stage,
	}
}

