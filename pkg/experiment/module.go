package experiment


// verify module implements Module interface
var _ Module = (*module)(nil)

type module struct {
	shape Shape
	inPorts InPorts
	outPorts OutPorts
	subModules [] Module
	isCopied bool
	isFused bool
}

func (m *module) IsFused() bool {
	return m.isFused
}

func (m *module) Shape() Shape {
	return m.shape
}

func (m *module) ReplaceShape(shape Shape) Module {
	panic("implement me")
}

func (m *module) InPorts() []InPort {
	return m.inPorts
}

func (m *module) OutPorts() []OutPort {
	return m.outPorts
}

func (m *module) IsRunnable() bool {
	return len(m.inPorts) ==0 && len(m.outPorts) ==0
}

func (m *module) IsSink() bool {
	return len(m.inPorts) ==1 && len(m.outPorts) ==0
}

func (m *module) IsSource() bool {
	return len(m.inPorts) ==0 && len(m.outPorts) ==1
}

func (m *module) IsFlow() bool {
	return len(m.inPorts) ==1 && len(m.outPorts) ==1
}

func (m *module) IsBidiFlow() bool {
	return len(m.inPorts) ==2 && len(m.outPorts) ==2
}

func (m *module) IsAtomic() bool {
	return len(m.subModules) > 0
}

func (m *module) IsCopied() bool {
	return m.isCopied
}

func (m *module) Fuse(that Module, from OutPort, to InPort, f MatFunc) Module {
	return m.Compose(that, f).Wire(from, to)
}

func (m *module) Wire(from OutPort, to InPort) Module {
	if !m.outPorts.contains(from) {

	}

	if !m.inPorts.contains(to) {

	}

	return m
}

func (m *module) TransformMaterializedValue(f func(value interface{}) interface{}) Module {
	panic("implement me")
}

func (m *module) Compose(other Module, f MatFunc) Module {
	panic("implement me")
}

func (m *module) ComposeNoMaterialized(that Module) Module {
	panic("implement me")
}

func (m *module) Nest() Module {
	panic("implement me")
}

func (m *module) SubModules() []Module {
	panic("implement me")
}

func (m *module) IsSealed() bool {
	return m.IsAtomic() || m.IsCopied() || m.IsFused()
}

func (m *module) Downstreams() map[OutPort]InPort {
	panic("implement me")
}

func (m *module) Upstreams() map[InPort]OutPort {
	panic("implement me")
}

func (m *module) MaterializedValueComputation() MaterializedValueNode {
	panic("implement me")
}

func (m *module) CarbonCopy() Module {
	panic("implement me")
}

func (m *module) Attributes() Attributes {
	panic("implement me")
}

func (m *module) WithAttributes(attributes Attributes) Module {
	panic("implement me")
}

