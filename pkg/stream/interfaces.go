package stream

import "context"

type Stage interface {
	Name() string
	Run(ctx context.Context)
}

type SinkStage interface {
	InputStage
	Result() Future
}
type SinkStageWithOptions interface {
	SinkStage
	WithOptions(options ...StageOption) SinkStageWithOptions
}

type Future interface {
	Await() (interface{}, error)
}

type OutputStage interface {
	Stage
	Outlet() *Outlet
}

type InputStage interface {
	Stage
	Inlet() * Inlet
}

type SourceStage interface {
	OutputStage
}

type SourceStageWithOptions interface {
	SourceStage
	WithOptions(options ...StageOption) SourceStageWithOptions
}

type FlowStage interface {
	InputStage
	OutputStage
}

type FlowStageWithOptions interface {
	InputStage
	SourceStage
	WithOptions(options ...StageOption) FlowStageWithOptions
}