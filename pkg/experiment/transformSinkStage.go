package experiment

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
)

// verify transformSinkStage implements SinkStage interface
var _ SinkStage = (*transformSinkStage)(nil)

type transformSinkStage struct {
	sinkStage SinkStage
	f MapMaterializedValueFunc
}

func (t *transformSinkStage) With(options ...StageOption) Stage {
	return &transformSinkStage{
		sinkStage: t.sinkStage.With(options...).(SinkStage),
		f: t.f,
	}
}

func (t *transformSinkStage) WireTo(stage UpstreamStage) SinkStage {
	return &transformSinkStage{
		sinkStage: t.sinkStage.WireTo(stage),
		f: t.f,
	}
}

func (t *transformSinkStage) Run(ctx context.Context, mat MaterializeFunc) *core.Promise {
	return t.Run(ctx, mat).FlatMap(t.f)
}

func TransformSink(sink SinkStage, f MapMaterializedValueFunc) SinkStage {
	return &transformSinkStage{
		sinkStage: sink,
		f:         f,
	}
}

