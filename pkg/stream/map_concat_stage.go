package stream

import (
	"fmt"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/core/iterable"
)

type MapConcatFunc func(value interface{}) (iterable.Iterable, error)

// verify mapConcatFlowLogic implements FlowStageLogic interface
var _ FlowStageLogic = (*mapConcatFlowLogic)(nil)

type mapConcatFlowLogic struct {
	decider Decider
	f       MapConcatFunc
}

func (m *mapConcatFlowLogic) SupportsParallelism() bool {
	return true
}

func (m *mapConcatFlowLogic) OnUpstreamStart(_ FlowStageActions) {
}

func (m *mapConcatFlowLogic) OnUpstreamReceive(element Element, actions FlowStageActions) {
	element.
		WhenError(actions.SendError).
		WhenValue(m.handleValue(actions))
}

func (m *mapConcatFlowLogic) OnUpstreamFinish(_ FlowStageActions) {
}

func (m *mapConcatFlowLogic) handleValue(actions FlowStageActions) ValueAction {
	return func(value interface{}) {
		iterable, err := m.f(value)
		if err != nil {
			switch m.decider(err) {
			case Stop:
				actions.FailStage(err)
			case Resume:
				actions.SendError(err)
			case Reset:
			}
			return
		}
		iterable.TakeWhile(func(value core.Any) bool {
			return !actions.StageIsCompleted()
		}).ForEach(func(item interface{}) {
			if actions.StageIsCompleted() {
				fmt.Println("Stage is complete")
				return
			}
			actions.SendValue(item)
		})
	}
}

func mapConcatStage(f MapConcatFunc) FlowStage {
	return Flow(mapConcatFactory(f)).Named("MapConcatStage").(FlowStage)
}

func mapConcatFactory(f MapConcatFunc) FlowStageLogicFactory {
	return func(attributes *StageAttributes) FlowStageLogic {
		return &mapConcatFlowLogic{
			decider: attributes.Decider,
			f:       f,
		}
	}
}
