package flow

import (
	"sync"

	"github.com/nicholasham/piper/pkg/zz/stream"
)

// verify scanOperator implements OperatorLogic interface
var _ OperatorLogic = (*scanOperator)(nil)

type scanOperator struct {
	current interface{}
	f       AggregateFunc
	sync.RWMutex
}

func (receiver *scanOperator) SupportsParallelism() bool {
	return false
}

func (receiver *scanOperator) Start(actions OperatorActions) {
	actions.PushValue(receiver.current)
}

func (receiver *scanOperator) Apply(element stream.Element, actions OperatorActions) {
	receiver.Lock()
	defer receiver.Unlock()
	out, err := receiver.f(receiver.current, element.Value())
	if err != nil {
		actions.FailStage(err)
	}
	receiver.current = out
	actions.PushValue(out)
}

func (receiver *scanOperator) End(actions OperatorActions) {
}

func scan(zero interface{}, f AggregateFunc) OperatorLogic {
	return &scanOperator{
		current: zero,
		f:       f,
	}
}
