package flow

type SinkGraph struct {
	stage SinkStage
}

func (g *SinkGraph) WithAttributes(attributes ...Attribute) *SinkGraph {
	return FromSink(g.stage.WithAttributes(attributes...))
}

func FromSink(stage SinkStage) *SinkGraph {
	return &SinkGraph{
		stage: stage,
	}
}
