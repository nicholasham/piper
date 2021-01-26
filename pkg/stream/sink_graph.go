package stream

type SinkGraph struct {
	stage  SinkStage
}

func (g *SinkGraph) Stage() SinkStage {
	return g.stage
}

func SinkFrom(sinkStage SinkStage) *SinkGraph {
	return &SinkGraph{
		stage:  sinkStage,
	}
}
