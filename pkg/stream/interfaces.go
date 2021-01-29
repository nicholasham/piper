package stream

import "context"

type Stage interface {
	Name() string
	Run(ctx context.Context)
	With(options ...StageOption) Stage

}

type SinkStage interface {
	InputStage
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
	WireTo(stage OutputStage)
}

type SourceStage interface {
	OutputStage
}

type FlowStage interface {
	InputStage
	OutputStage
}

