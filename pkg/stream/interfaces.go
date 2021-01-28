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
	With(options ...StageOption) SinkStageWithOptions
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
	WireTo(stage OutputStage)
}

type SourceStage interface {
	OutputStage
}

type SourceStageWithOptions interface {
	SourceStage
	With(options ...StageOption) SourceStageWithOptions
}

type FlowStage interface {
	InputStage
	OutputStage
}

type FlowStageWithOptions interface {
	InputStage
	SourceStage
	With(options ...StageOption) FlowStageWithOptions
}