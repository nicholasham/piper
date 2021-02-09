package experiment

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

func (g *SourceGraph) viaFlow(that FlowStage) *SourceGraph  {
	return FromSource(that.WireTo(g.stage))
}

func (g *SourceGraph) To(that *SinkGraph) *RunnableGraph {
	return g.ToMaterialized(that)(KeepLeft)
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

func (g *SourceGraph) MapConcat(f MapConcatFunc) *SourceGraph {
	return g.viaFlow(MapConcat(f))
}

func FromSource(stage SourceStage) *SourceGraph {
	return &SourceGraph{
		stage: stage,
	}
}


