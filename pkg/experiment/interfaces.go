package experiment


type Attributes struct {

}

type InPort interface {

}

type OutPort interface {

}

type MatFunc  func(a interface{}, b interface{}) interface{}


type Module interface {
	Shape() Shape
	ReplaceShape(shape Shape) Module
	InPorts() [] InPort
	OutPorts() [] OutPort
	IsRunnable() bool
	IsSink() bool
	IsSource() bool
	IsFlow() bool
	IsBidiFlow() bool
	IsAtomic() bool
	IsCopied() bool
	IsFused() bool
	Fuse(that Module, from OutPort, to InPort, f MatFunc) Module
	Wire(from OutPort, to InPort)
	TransformMaterializedValue(f func(value interface{}) interface{}) Module
	Compose(that Module) Module
	ComposeNoMaterialized(that Module) Module
	Nest() Module
	SubModules()[]Module
	IsSealed() bool
	Downstreams() map[OutPort]InPort
	Upstreams() map[InPort]OutPort
	MaterializedValueComputation() MaterializedValueNode
	CarbonCopy() Module
	Attributes() Attributes
	WithAttributes(attributes Attributes) Module
}

type MaterializedValueNode interface {

}

type Shape interface {

}

type Graph interface {
	Module() Module
	WithAttributes(attributes Attributes) Graph
	AddAttributes(attributes Attributes) Graph
	Named(name string) Graph
}





