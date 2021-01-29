package stream

import "context"

type Stage interface {
	Name() string
	Run(ctx context.Context)
	With(options ...StageOption) Stage
}

type SinkStage interface {
	InputStage
	WireTo(stage OutputStage) SinkStage
	Result() Future
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

}

type SourceStage interface {
	OutputStage
}

type FlowStage interface {
	InputStage
	OutputStage
	WireTo(stage OutputStage) FlowStage
}
