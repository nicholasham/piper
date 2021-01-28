package stream

type SourceGraph struct {
	stage SourceStage
}

func (g *SourceGraph) WithOptions(options ...StageOption) *SourceGraph {
	return FromSource(g.stage.WithOptions(options...))
}

func FromSource(stage SourceStage) *SourceGraph {
	return &SourceGraph{
		stage: stage,
	}
}

