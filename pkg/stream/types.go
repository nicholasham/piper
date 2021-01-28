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

type OutStage interface {
	Stage
	Wire(stage SourceStage)
	Outlet() *Outlet
}

type SourceStage interface {
	OutStage
	WithOptions(options ...StageOption) SourceStage
}

type FlowStage interface {
	OutStage
	WithOptions(options ...StageOption) FlowStage
}

type CompletionStage interface {
	Result() Future
}

