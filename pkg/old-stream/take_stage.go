package old_stream

import "sync/atomic"

// verify takeFlowLogic implements FlowStageLogic interface
var _ FlowStageLogic = (*takeFlowLogic)(nil)

type takeFlowLogic struct {
	op     uint64
	number int
}

func (t *takeFlowLogic) SupportsParallelism() bool {
	return true
}

func (t *takeFlowLogic) OnUpstreamStart(actions FlowStageActions) {
}

func (t *takeFlowLogic) OnUpstreamReceive(element Element, actions FlowStageActions) {
	element.
		WhenError(actions.SendError).
		WhenValue(t.handleValue(actions))
}

func (t *takeFlowLogic) handleValue(actions FlowStageActions) ValueAction {
	return func(value interface{}) {
		current := int(atomic.AddUint64(&t.op, 1))
		if current <= t.number {
			actions.SendValue(value)
		}
	}
}

func (t *takeFlowLogic) OnUpstreamFinish(actions FlowStageActions) {
}

func Take(number int) FlowStage {
	return LinearFlow(takeFactory(number))
}

func takeFactory(number int) FlowStageLogicFactory {
	return func(attributes *StageAttributes) FlowStageLogic {
		return &takeFlowLogic{
			number: number,
		}
	}
}
