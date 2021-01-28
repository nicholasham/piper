package stream

type SinkGraph struct {
	stage SinkStageWithOptions
}

func (g *SinkGraph) WithOptions(options ...StageOption) *SinkGraph {
	return FromSink(g.stage.WithOptions(options...))
}

func FromSink(stage SinkStageWithOptions) *SinkGraph {
	return &SinkGraph{
		stage: stage,
	}
}
