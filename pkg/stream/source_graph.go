package stream

type SourceGraph struct {
	stage SourceStage
}

func (g *SourceGraph) With(options ...StageOption) *SourceGraph {
	return FromSource(g.stage.With(options...).(SourceStage))
}

func (g *SourceGraph) Named(name string) *SourceGraph {
	return FromSource(g.stage.With(Name(name)).(SourceStage))
}

func (g *SourceGraph) Via(that *FlowGraph) *FlowGraph {
	return FromFlow(NewFusedFlow(g.stage, that.stage))
}

func FromSource(stage SourceStage) *SourceGraph {
	return &SourceGraph{
		stage: stage,
	}
}

