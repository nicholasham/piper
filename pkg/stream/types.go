package stream

import "context"

type Stage interface {
	Name() string
	Run(ctx context.Context)
}


type SinkStage interface {
	Stage
	Wire(stage SourceStage)
	Inlet() *Inlet
	Result() Future
	WithOptions(options ...StageOption) SinkStage
}

type Future interface {
	Await() (interface{}, error)
}

type SourceStage interface {
	Stage
	Outlet() *Outlet
	WithOptions(options ...StageOption) SourceStage
}

type FlowStage interface {
	Wire(stage SourceStage)
	WithOptions(options ...StageOption) FlowStage
}


type CompletionStage interface {
	Result() Future
}

