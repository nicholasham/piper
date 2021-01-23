package stream

import (
	"context"
)

type RunnableGraph struct {
	stages []Stage
	sink   SinkStage
}

func (r *RunnableGraph) Run(ctx context.Context) Future {
	for _, stage := range r.stages {
		stage.Run(ctx)
	}
	return r.sink.Result()
}

func sourceRunnable(sourceGraph *SourceGraph, sinkGraph *SinkGraph) *RunnableGraph {
	sinkGraph.stage.WireTo(sourceGraph.stage)
	combinedStages := combineStages(sourceGraph.stages, sinkGraph.stages)

	return &RunnableGraph{
		stages: combinedStages,
		sink:   sinkGraph.stage,
	}
}

func flowRunnable(flowGraph *FlowGraph, sinkGraph *SinkGraph) *RunnableGraph {
	sinkGraph.stage.WireTo(flowGraph.stage)
	combinedStages := combineStages(flowGraph.stages, sinkGraph.stages)
	return &RunnableGraph{
		stages: combinedStages,
		sink:   sinkGraph.stage,
	}
}
