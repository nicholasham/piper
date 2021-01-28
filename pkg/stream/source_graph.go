package stream

type SourceGraph struct {
	stage SourceStageWithOptions
}

func (g *SourceGraph) With(options ...StageOption) *SourceGraph {
	return FromSource(g.stage.With(options...))
}

func (g *SourceGraph) Via(that *FlowGraph) *FlowGraph {
	return FromFlow(NewFusedFlow(g.stage, that.stage))
}

func FromSource(stage SourceStageWithOptions) *SourceGraph {
	return &SourceGraph{
		stage: stage,
	}
}

