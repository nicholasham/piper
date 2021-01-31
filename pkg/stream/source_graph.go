package stream

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
	return FromFlow(NewFusedFlow(g.stage, that.stage))
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
	return FromSource(Map(f).WireTo(g.stage))
}

func (g *SourceGraph) MapConcat(f MapConcatFunc) *SourceGraph {
	return FromSource(MapConcat(f).WireTo(g.stage))
}

func (g *SourceGraph) Filter(f FilterFunc) *SourceGraph {
	return FromSource(Filter(f).WireTo(g.stage))
}

func (g *SourceGraph) Drop(number int) *SourceGraph {
	return FromSource(Drop(number).WireTo(g.stage))
}

func (g *SourceGraph) Take(number int) *SourceGraph {
	return FromSource(Take(number).WireTo(g.stage))
}

func (g *SourceGraph) TakeWhile(f FilterFunc) *SourceGraph {
	return FromSource(TakeWhile(f).WireTo(g.stage))
}

func (g *SourceGraph) Fold(zero interface{}, f AggregateFunc) *SourceGraph {
	return FromSource(Fold(zero, f).WireTo(g.stage))
}

func (g *SourceGraph) Unfold(state interface{}, f UnfoldFunc) *SourceGraph {
	return FromSource(Unfold(state, f).WireTo(g.stage))
}


