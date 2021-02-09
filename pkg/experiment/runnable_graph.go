package experiment

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
)

type RunnableGraph struct {
	combine MaterializeFunc
	sinkStage   SinkStage
}

func (r *RunnableGraph) Run(ctx context.Context) *core.Future {
	return r.sinkStage.Run(ctx, r.combine)
}

func runnable(sinkStage SinkStage, combine MaterializeFunc) *RunnableGraph {
	return &RunnableGraph{
		sinkStage:   sinkStage,
		combine: combine,
	}
}
