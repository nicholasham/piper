package flow

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
	WithAttributes(attributes ...Attribute) SinkStage
}

type Future interface {
	Await() (interface{}, error)
}

type SourceStage interface {
	Stage
	Outlet() *Outlet
	WithAttributes(attributes ...Attribute) SourceStage
}

type FlowStage interface {
	Wire(stage SourceStage)
	WithAttributes(attributes ...Attribute) FlowStage
}


type CompletionStage interface {
	Result() Future
}

