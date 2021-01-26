package stream

import (
	"context"
)

type RunnableGraph struct {
	sourceStage SourceStage
	sinkStage   SinkStage
}

func (r *RunnableGraph) Run(ctx context.Context) Future {
	r.sourceStage.Run(ctx)
	r.sinkStage.Run(ctx)
	return r.sinkStage.Result()
}

func runnable(sourceStage SourceStage,  sinkStage SinkStage ) *RunnableGraph {
	sinkStage.Wire(sourceStage)
	return &RunnableGraph{
		sourceStage: sourceStage,
		sinkStage:   sinkStage,
	}
}


