package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
)

type SourceGraph struct {
	stage SourceStage
}

func (g *SourceGraph) With(options ...StageOption) *SourceGraph {
	return FromSource(g.stage.With(options...).(SourceStage))
}

func (g *SourceGraph) Named(name string) *SourceGraph {
	return g.With(Name(name))
}

func (g *SourceGraph) Via(that *FlowGraph) *FlowGraph {
	that.stage.WireTo(g.stage)
	return that
}

func (g *SourceGraph) viaFlow(that FlowStage) *SourceGraph {
	return FromSource(that.WireTo(g.stage))
}

func (g *SourceGraph) To(that *SinkGraph) *RunnableGraph {
	return g.ToMaterialized(that)(KeepRight)
}

func (g *SourceGraph) AlsoTo(that *SinkGraph) *SourceGraph {
	return g.viaFlow(diversion(that.stage, alsoToStrategy()))
}

func (g *SourceGraph) DivertTo(that *SinkGraph, when core.PredicateFunc) *SourceGraph {
	return g.viaFlow(diversion(that.stage, divertToStrategy(when)))
}

func (g *SourceGraph) ToMaterialized(that *SinkGraph) func(combine MaterializeFunc) *RunnableGraph {
	return func(combine MaterializeFunc) *RunnableGraph {
		that.stage.WireTo(g.stage)
		return runnable(that.stage, combine)
	}
}

func (g *SourceGraph) RunWith(ctx context.Context, that *SinkGraph) *core.Future {
	return g.ToMaterialized(that)(KeepRight).Run(ctx)
}

func (g *SourceGraph) Drop(number int) *SourceGraph {
	return g.viaFlow(dropStage(number))
}

func (g *SourceGraph) Filter(f FilterFunc) *SourceGraph {
	return g.viaFlow(filterStage(f))
}

func (g *SourceGraph) Fold(zero interface{}, f AggregateFunc) *SourceGraph {
	return g.viaFlow(foldStage(zero, f))
}

func (g *SourceGraph) Map(f MapFunc) *SourceGraph {
	return g.viaFlow(mapStage(f))
}

func (g *SourceGraph) MapConcat(f MapConcatFunc) *SourceGraph {
	return g.viaFlow(mapConcatStage(f))
}

func (g *SourceGraph) Scan(zero interface{}, f AggregateFunc) *SourceGraph {
	return g.viaFlow(scanStage(zero, f))
}

func (g *SourceGraph) Take(number int) *SourceGraph {
	return g.viaFlow(takeStage(number))
}

func (g *SourceGraph) TakeWhile(f FilterFunc) *SourceGraph {
	return g.viaFlow(takeWhileStage(f))
}

func (g *SourceGraph) Unfold(state interface{}, f UnfoldFunc) *SourceGraph {
	return g.viaFlow(unfoldStage(state, f))
}

func FromSource(stage SourceStage) *SourceGraph {
	return &SourceGraph{
		stage: stage,
	}
}
