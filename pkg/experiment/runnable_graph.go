package experiment

import (
	"context"
)

type RunnableGraph struct {
	combine MaterializeFunc
	sinkStage   SinkStage
}

func (r *RunnableGraph) Run(ctx context.Context) Future {
	return r.sinkStage.Run(ctx, r.combine)
}

func runnable(sinkStage SinkStage, combine MaterializeFunc) *RunnableGraph {
	return &RunnableGraph{
		sinkStage:   sinkStage,
		combine: combine,
	}
}
