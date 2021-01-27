package flow

type SourceGraph struct {
	stage SourceStage
}

func (g *SourceGraph) WithAttributes(attributes ...Attribute) *SourceGraph {
	return FromSource(g.stage.WithAttributes(attributes...))
}

func FromSource(stage SourceStage) *SourceGraph {
	return &SourceGraph{
		stage: stage,
	}
}

