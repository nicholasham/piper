package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
)

type Stage interface {
	With(options ...StageOption) Stage
	Named(name string) Stage
}

type SourceStage interface {
	Stage
	UpstreamStage
}

type FlowStage interface {
	Stage
	UpstreamStage
	WireTo(stage UpstreamStage) FlowStage
}

type SinkStage interface {
	Stage
	WireTo(stage UpstreamStage) SinkStage
	Run(ctx context.Context, mat MaterializeFunc) *core.Future
}

type UpstreamStage interface {
	Open(ctx context.Context, mat MaterializeFunc) (Reader, *core.Future)
}
