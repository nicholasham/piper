package piper

import "context"

type Stage interface {
	Run(ctx context.Context)
}


type Shape interface {
	Inlets() []Inlet
	Outlets() []Outlet
}

type FlowShape interface {
	Inlet() Inlet
	Outlet() Outlet
}

type Flow interface {
	Via(flow Flow) Flow
}

type SourceGraph interface {
	Shape() SourceShape
}

type FlowGraph interface {
	 Shape() FlowShape
}

type SinkGraph interface {
	Shape() SinkShape
	RunWith(source SourceGraph) Future
}

type RunnableGraph interface {
	Shape() ClosedShape
	Run() Future
}

type ClosedShape interface {

}

type SourceShape interface {
	Inlet() Inlet
}

type SinkShape interface {
	Inlet() Inlet
}

type FanInShape interface {
	Inlets() [] Inlet
	Outlet() [] Outlet
}

type FanOutShape interface {
	Inlet() Inlet
	Outlets() Outlet
}



type Future interface {
	Await() (interface{}, error)
}


