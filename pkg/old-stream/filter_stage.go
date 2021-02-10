package old_stream

type FilterFunc func(value interface{}) bool

// verify filterFlowLogic implements FlowStageLogic interface
var _ FlowStageLogic = (*filterFlowLogic)(nil)

type filterFlowLogic struct {
	f FilterFunc
}

func (f *filterFlowLogic) SupportsParallelism() bool {
	return true
}

func (f *filterFlowLogic) OnUpstreamStart(actions FlowStageActions) {
}

func (f *filterFlowLogic) OnUpstreamReceive(element Element, actions FlowStageActions) {
	element.
		WhenError(actions.SendError).
		WhenValue(f.handleValue(actions))
}

func (f *filterFlowLogic) OnUpstreamFinish(actions FlowStageActions) {
}

func (f *filterFlowLogic) handleValue(actions FlowStageActions) ValueAction {
	return func(value interface{}) {
		if f.f(value) {
			actions.SendValue(value)
		}
	}
}

func Filter(f FilterFunc) FlowStage {
	return LinearFlow(filterFactory(f))
}

func filterFactory(f FilterFunc) FlowStageLogicFactory {
	return func(attributes *StageAttributes) FlowStageLogic {
		return &filterFlowLogic{
			f: f,
		}
	}
}
