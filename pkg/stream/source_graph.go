package stream

type SourceGraph struct {
	stage SourceStage
}

func (g *SourceGraph) WithOptions(options ...StageOption) *SourceGraph {
	return FromSource(g.stage.WithOptions(options...))
}

func (g *SourceGraph) Via(that *FlowGraph) *FlowGraph {
	return FromFlow(NewFusedFlow(g.stage, that.stage))
}

func FromSource(stage SourceStage) *SourceGraph {
	return &SourceGraph{
		stage: stage,
	}
}

