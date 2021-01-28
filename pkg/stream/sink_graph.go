package stream

type SinkGraph struct {
	stage SinkStageWithOptions
}

func (g *SinkGraph) With(options ...StageOption) *SinkGraph {
	return FromSink(g.stage.With(options...))
}

func FromSink(stage SinkStageWithOptions) *SinkGraph {
	return &SinkGraph{
		stage: stage,
	}
}
