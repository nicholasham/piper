package old_stream

import "context"

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
	return FromFlow(CompositeFlow(g.stage, that.stage))
}

func (g *SourceGraph) viaFlow(that FlowStage) *SourceGraph {
	return FromSource(CompositeFlow(g.stage, that))
}

func (g *SourceGraph) To(that *SinkGraph) *RunnableGraph {
	return runnable(g.stage, that.stage)
}

func (g *SourceGraph) RunWith(ctx context.Context, that *SinkGraph) Future {
	return g.To(that).Run(ctx)
}

func FromSource(stage SourceStage) *SourceGraph {
	return &SourceGraph{
		stage: stage,
	}
}

func (g *SourceGraph) Concat(that *SourceGraph) *SourceGraph {
	return ConcatSources(g, that)
}

func (g *SourceGraph) Interleave(segmentSize int, that *SourceGraph) *SourceGraph {
	return InterleaveSources(segmentSize, g, that)
}

func (g *SourceGraph) Merge(that *SourceGraph) *SourceGraph {
	return MergeSources(g, that)
}

func (g *SourceGraph) Map(f MapFunc) *SourceGraph {
	return g.viaFlow(Map(f))
}

func (g *SourceGraph) MapConcat(f MapConcatFunc) *SourceGraph {
	return g.viaFlow(MapConcat(f))
}

func (g *SourceGraph) Filter(f FilterFunc) *SourceGraph {
	return g.viaFlow(Filter(f))
}

func (g *SourceGraph) Drop(number int) *SourceGraph {
	return g.viaFlow(Drop(number))
}

func (g *SourceGraph) Take(number int) *SourceGraph {
	return g.viaFlow(Take(number))
}

func (g *SourceGraph) TakeWhile(f FilterFunc) *SourceGraph {
	return g.viaFlow(TakeWhile(f))
}

func (g *SourceGraph) Fold(zero interface{}, f AggregateFunc) *SourceGraph {
	return g.viaFlow(Fold(zero, f))
}

func (g *SourceGraph) Unfold(state interface{}, f UnfoldFunc) *SourceGraph {
	return g.viaFlow(Unfold(state, f))
}
