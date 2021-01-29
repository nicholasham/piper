package stream

import "context"

type Stage interface {
	Name() string
	Run(ctx context.Context)
}

type InputStage interface {
	Inlet() *Inlet
	Wire(stage SourceStage)
	Stage
}

type OutputStage interface {
	Outlet() *Outlet
	Stage
}

type SinkStage interface {
	Stage
	Wire(stage SourceStage)
	Inlet() *Inlet
	Result() Future
}

type Future interface {
	Await() (interface{}, error)
}

type SourceStage interface {
	Stage
	Outlet() *Outlet
}

type FlowStage interface {
	SourceStage
	Wire(stage SourceStage)
}

type CompletionStage interface {
	Result() Future
}
