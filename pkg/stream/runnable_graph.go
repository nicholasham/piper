package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
	"sync"
)

type RunnableGraph struct {
	combine   MaterializeFunc
	sinkStage SinkStage
}

func (r *RunnableGraph) Run(ctx context.Context) *core.Future {
	wg := &sync.WaitGroup{}
	future:= r.sinkStage.Run(ctx, wg, r.combine)
	return core.NewFuture(func() core.Result {
		wg.Wait()
		return future.Await()
	})
}

func runnable(sinkStage SinkStage, combine MaterializeFunc) *RunnableGraph {
	return &RunnableGraph{
		sinkStage: sinkStage,
		combine:   combine,
	}
}
