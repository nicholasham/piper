package experiment

import (
	"context"
)

type RunnableGraph struct {
	sinkStage   SinkStage
}

func (r *RunnableGraph) Run(ctx context.Context, mat MaterializeFunc) Future {
	return r.sinkStage.Run(ctx, mat)
}

func runnable(sourceStage SourceStage, sinkStage SinkStage) *RunnableGraph {
	sinkStage.WireTo(sourceStage)
	return &RunnableGraph{
		sinkStage:   sinkStage,
	}
}
