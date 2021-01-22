package stream

import "context"

type RunSink func(ctx context.Context, inlet Inlet, result Promise)
type RunSource func(ctx context.Context, outlet Outlet)
type RunFlow func(ctx context.Context, inlet Inlet, outlet Outlet)

type Stage interface {
	Name() string
	Run(ctx context.Context)
}

type SinkStage interface {
	Stage
	WireTo(stage SourceStage)
	Inlet() *Inlet
	Result() Future
}

type Future interface {
	Await() (interface{}, error)
}

type SourceStageFactory func(options ...Option) SourceStage

type SourceStage interface {
	Stage
	Outlet() *Outlet
}

type FlowStage interface {
	SourceStage
	Wire(stage SourceStage)
}
