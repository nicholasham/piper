package stream

type MapFunc func(value interface{}) (interface{}, error)

// verify mapFlowLogic implements FlowStageLogic interface
var _ FlowStageLogic = (*mapFlowLogic)(nil)

type mapFlowLogic struct {
	decider Decider
	f       MapFunc
}

func (f *mapFlowLogic) SupportsParallelism() bool {
	return true
}

func (f *mapFlowLogic) OnUpstreamStart(actions FlowStageActions) {
}

func (f *mapFlowLogic) OnUpstreamReceive(element Element, actions FlowStageActions) {
	element.
		WhenError(actions.SendError).
		WhenValue(f.handleValue(actions))
}

func (f *mapFlowLogic) OnUpstreamFinish(actions FlowStageActions) {
}

func (f *mapFlowLogic) handleValue(actions FlowStageActions) ValueAction {
	return func(fromValue interface{}) {
		toValue, err := f.f(fromValue)
		if err != nil {
			switch f.decider(err) {
			case Stop:
				actions.FailStage(err)
			case Resume:
				actions.SendError(err)
			case Reset:
			}
			return
		}
		actions.SendValue(toValue)
	}
}

func mapStage(f MapFunc) FlowStage {
	return Flow(mapFactory(f))
}

func mapFactory(f MapFunc) FlowStageLogicFactory {
	return func(attributes *StageAttributes) FlowStageLogic {
		return &mapFlowLogic{
			f: f,
		}
	}
}
