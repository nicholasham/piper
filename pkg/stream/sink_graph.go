package stream

type SinkGraph struct {
	stage  SinkStage
	stages []Stage
}

func (g *SinkGraph) Stage() SinkStage {
	return g.stage
}

func SinkFrom(sinkStage SinkStage, upstreamStages ...Stage) *SinkGraph {
	return &SinkGraph{
		stages: append(upstreamStages, sinkStage),
		stage:  sinkStage,
	}
}
