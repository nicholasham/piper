package piper

import "context"

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

type SourceStage interface {
	Stage
	Outlet() *Outlet
}

type FlowStage interface {
	SourceStage
	Wire(stage SourceStage)
}
