package old_stream

// verify takeWhileFlowLogic implements FlowStageLogic interface
var _ FlowStageLogic = (*takeWhileFlowLogic)(nil)

type takeWhileFlowLogic struct {
	f FilterFunc
}

func (t *takeWhileFlowLogic) SupportsParallelism() bool {
	return true
}

func (t *takeWhileFlowLogic) OnUpstreamStart(actions FlowStageActions) {
}

func (t *takeWhileFlowLogic) OnUpstreamReceive(element Element, actions FlowStageActions) {
	element.
		WhenError(actions.SendError).
		WhenValue(t.handleValue(actions))
}

func (t *takeWhileFlowLogic) handleValue(actions FlowStageActions) ValueAction {
	return func(value interface{}) {
		if t.f(value) {
			actions.SendValue(value)
		}
	}
}

func (t *takeWhileFlowLogic) OnUpstreamFinish(actions FlowStageActions) {
}

func TakeWhile(f FilterFunc) FlowStage {
	return LinearFlow(takeWhileFactory(f))
}

func takeWhileFactory(f FilterFunc) FlowStageLogicFactory {
	return func(attributes *StageAttributes) FlowStageLogic {
		return &takeWhileFlowLogic{
			f: f,
		}
	}
}
