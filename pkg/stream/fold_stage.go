package stream

import "sync"

type AggregateFunc func(acc interface{}, out interface{}) (interface{}, error)

// verify foldFlowStage implements FlowStageLogic interface
var _ FlowStageLogic = (*foldFlowStage)(nil)

type foldFlowStage struct {
	current interface{}
	f       AggregateFunc
	sync.RWMutex
	decider Decider
}

func (f *foldFlowStage) SupportsParallelism() bool {
	return true
}

func (f *foldFlowStage) OnUpstreamStart(actions FlowStageActions) {
}

func (f *foldFlowStage) OnUpstreamReceive(element Element, actions FlowStageActions) {
}

func (f *foldFlowStage) OnUpstreamFinish(actions FlowStageActions) {
	actions.SendValue(f.current)
}

func (f *foldFlowStage) handleValue(actions FlowStageActions) ValueAction {
	return func(value interface{}) {
		f.Lock()
		defer f.Unlock()
		out, err := f.f(f.current, value)
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
		f.current = out
	}
}

func Fold(zero interface{}, f AggregateFunc) FlowStage {
	return LinearFlow(foldFactory(zero, f))
}

func foldFactory(zero interface{}, f AggregateFunc) FlowStageLogicFactory {
	return func(attributes *StageAttributes) FlowStageLogic {
		return &foldFlowStage{
			current: zero,
			f:       f,
			decider: attributes.Decider,
		}
	}
}
