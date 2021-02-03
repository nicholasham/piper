package experiment


type Attributes struct {

}

type InPort interface {

}

type OutPort interface {

}

type InPorts [] InPort
type OutPorts [] OutPort

func (i InPorts) contains(inPort InPort) bool {
	for _, p := range i {
		if p== inPort {
			return true
		}
	}
	return false
}

func (o OutPorts) contains(outPort InPort) bool {
	for _, p := range o {
		if p == outPort {
			return true
		}
	}
	return false
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
	Wire(from OutPort, to InPort) Module
	TransformMaterializedValue(f func(value interface{}) interface{}) Module
	Compose(that Module, f MatFunc) Module
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





