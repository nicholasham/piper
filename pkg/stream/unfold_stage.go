package stream

import (
	"github.com/nicholasham/piper/pkg/core"
	"sync"
)

type UnfoldFunc func(state interface{}) core.Optional

// verify unfoldFlowStage implements FlowStageLogic interface
var _ FlowStageLogic = (*unfoldFlowStage)(nil)

type unfoldFlowStage struct {
	result core.Optional
	f      UnfoldFunc
	sync.RWMutex
}

func (u *unfoldFlowStage) SupportsParallelism() bool {
	return true
}

func (u *unfoldFlowStage) OnUpstreamStart(actions FlowStageActions) {
}

func (u *unfoldFlowStage) OnUpstreamReceive(element Element, actions FlowStageActions) {
	element.
		WhenError(actions.SendError).
		WhenValue(u.handleValue(actions))
}

func (u *unfoldFlowStage) handleValue(actions FlowStageActions) ValueAction { 
	return func(value interface{}) {
		u.Lock()
		defer u.Unlock()

		u.result.
			IfSome(actions.SendValue).
			IfNone(actions.CompleteStage)

		u.result = u.f(value)
	}
}

func (u *unfoldFlowStage) OnUpstreamFinish(actions FlowStageActions) {
	panic("implement me")
}

func Unfold(state interface{}, f UnfoldFunc) FlowStage {
	return LinearFlow(unfoldFactory(state, f))
}

func unfoldFactory(state interface{}, f UnfoldFunc) FlowStageLogicFactory {
	return func(attributes *StageAttributes) FlowStageLogic {
		return & unfoldFlowStage{
			result: core.Some(state),
			f:      f,
		}
	}
}
