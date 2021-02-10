package old_stream

import "sync"

// verify scanFlowStage implements FlowStageLogic interface
var _ FlowStageLogic = (*scanFlowStage)(nil)

type scanFlowStage struct {
	current interface{}
	f       AggregateFunc
	sync.RWMutex
	decider Decider
}

func (s *scanFlowStage) SupportsParallelism() bool {
	return true
}

func (s *scanFlowStage) OnUpstreamStart(actions FlowStageActions) {
	actions.SendValue(s.current)
}

func (s *scanFlowStage) OnUpstreamReceive(element Element, actions FlowStageActions) {
	element.
		WhenError(actions.SendError).
		WhenValue(s.handleValue(actions))
}

func (s *scanFlowStage) OnUpstreamFinish(actions FlowStageActions) {

}

func (s *scanFlowStage) handleValue(actions FlowStageActions) ValueAction {
	return func(value interface{}) {
		s.Lock()
		defer s.Unlock()
		out, err := s.f(s.current, value)
		if err != nil {
			switch s.decider(err) {
			case Stop:
				actions.FailStage(err)
			case Resume:
				actions.SendError(err)
			case Reset:
			}
			return
		}
		s.current = out
		actions.SendValue(s.current)
	}
}

func Scan(zero interface{}, f AggregateFunc) FlowStage {
	return LinearFlow(scanFactory(zero, f))
}

func scanFactory(zero interface{}, f AggregateFunc) FlowStageLogicFactory {
	return func(attributes *StageAttributes) FlowStageLogic {
		return &scanFlowStage{
			current: zero,
			f:       f,
			decider: attributes.Decider,
		}
	}
}
