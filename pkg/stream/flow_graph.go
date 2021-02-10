package stream

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
	return g.viaFlow(mapStage(f))
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
