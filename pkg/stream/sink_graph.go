package stream

type SinkGraph struct {
	stage SinkStage
}

func (g *SinkGraph) With(options ...StageOption) *SinkGraph {
	return FromSink(g.stage.With(options...).(SinkStage))
}

func (g SinkGraph) MapMaterializedValue(f MapMaterializedValueFunc) *SinkGraph  {
	return FromSink(transformSink(g.stage, f))
}

func FromSink(stage SinkStage) *SinkGraph {
	return &SinkGraph{
		stage: stage,
	}
}
