package stream

import "sync/atomic"

// verify dropFlowLogic implements FlowStageLogic interface
var _ FlowStageLogic = (*dropFlowLogic)(nil)

type dropFlowLogic struct {
	op     uint64
	number int
}

func (d *dropFlowLogic) SupportsParallelism() bool {
	return true
}

func (d *dropFlowLogic) OnUpstreamStart(actions FlowStageActions) {
}

func (d *dropFlowLogic) OnUpstreamReceive(element Element, actions FlowStageActions) {
	element.
		WhenError(actions.SendError).
		WhenValue(d.handleValue(actions))
}

func (d *dropFlowLogic) handleValue(actions FlowStageActions) ValueAction {
	return func(value interface{}) {
		current := int(atomic.AddUint64(&d.op, 1))
		if current > d.number {
			actions.SendValue(value)
		}
	}
}

func Drop(number int) FlowStage {
	return LinearFlow(dropFactory(number))
}

func dropFactory(number int) FlowStageLogicFactory {
	return func(attributes *StageAttributes) FlowStageLogic {
		return &dropFlowLogic{
			number: number,
		}
	}
}

func (d *dropFlowLogic) OnUpstreamFinish(actions FlowStageActions) {
}
