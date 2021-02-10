package old_stream

type SinkGraph struct {
	stage SinkStage
}

func (g *SinkGraph) With(options ...StageOption) *SinkGraph {
	return FromSink(g.stage.With(options...).(SinkStage))
}

func FromSink(stage SinkStage) *SinkGraph {
	return &SinkGraph{
		stage: stage,
	}
}
