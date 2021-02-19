package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
)

type FlowGraph struct {
	stage FlowStage
}

func FromFlow(stage FlowStage) *FlowGraph {
	return &FlowGraph{
		stage: stage,
	}
}

func (g *FlowGraph) viaFlow(that FlowStage) *FlowGraph {
	return FromFlow(that.WireTo(g.stage))
}

func (g *FlowGraph) With(options ...StageOption) *FlowGraph {
	return FromFlow(g.stage.With(options...).(FlowStage))
}

func (g *FlowGraph) Named(name string) *FlowGraph {
	return g.With(Name(name))
}

func (g *FlowGraph) Via(that *FlowGraph) *FlowGraph {
	that.stage.WireTo(g.stage)
	return that
}

func (g *FlowGraph) To(that *SinkGraph) *RunnableGraph {
	return g.ToMaterialized(that)(KeepRight)
}

func (g *FlowGraph) ToMaterialized(that *SinkGraph) func(combine MaterializeFunc) *RunnableGraph {
	return func(combine MaterializeFunc) *RunnableGraph {
		that.stage.WireTo(g.stage)
		return runnable(that.stage, combine)
	}
}

func (g *FlowGraph) RunWith(ctx context.Context, that *SinkGraph) *core.Future {
	return g.ToMaterialized(that)(KeepRight).Run(ctx)
}

func (g *FlowGraph) AlsoTo(that *SinkGraph) *FlowGraph {
	return g.viaFlow(diversion(that.stage, alsoToStrategy()))
}


func (g *FlowGraph) DivertTo(that *SinkGraph, when core.PredicateFunc) *FlowGraph {
return g.viaFlow(diversion(that.stage, divertToStrategy(when)))
}

func (g *FlowGraph) Drop(number int) *FlowGraph {
	return g.viaFlow(dropStage(number))
}

func (g *FlowGraph) Filter(f FilterFunc) *FlowGraph {
	return g.viaFlow(filterStage(f))
}

func (g *FlowGraph) Fold(zero interface{}, f AggregateFunc) *FlowGraph {
	return g.viaFlow(foldStage(zero, f))
}

func (g *FlowGraph) Map(f MapFunc) *FlowGraph {
	return g.viaFlow(MapStage(f))
}

func (g *FlowGraph) MapConcat(f MapConcatFunc) *FlowGraph {
	return g.viaFlow(mapConcatStage(f))
}

func (g *FlowGraph) Scan(zero interface{}, f AggregateFunc) *FlowGraph {
	return g.viaFlow(scanStage(zero, f))
}

func (g *FlowGraph) Take(number int) *FlowGraph {
	return g.viaFlow(takeStage(number))
}

func (g *FlowGraph) TakeWhile(f FilterFunc) *FlowGraph {
	return g.viaFlow(takeWhileStage(f))
}

func (g *FlowGraph) Unfold(state interface{}, f UnfoldFunc) *FlowGraph {
	return g.viaFlow(unfoldStage(state, f))
}
