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
	return &RunnableGraph{
		stages: sourceGraph.combineStages(sinkGraph.stages),
		sink:   sinkGraph.stage,
	}
}

func flowRunnable(flowGraph *FlowGraph, sinkGraph *SinkGraph) *RunnableGraph {
	sinkGraph.stage.WireTo(flowGraph.stage)
	return &RunnableGraph{
		stages: flowGraph.combineStages(sinkGraph.stages),
		sink:   sinkGraph.stage,
	}
}
