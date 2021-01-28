package stream

type SinkGraph struct {
	stage SinkStage
}

func (g *SinkGraph) WithOptions(options ...StageOption) *SinkGraph {
	return FromSink(g.stage.WithOptions(options...))
}

func FromSink(stage SinkStage) *SinkGraph {
	return &SinkGraph{
		stage: stage,
	}
}
