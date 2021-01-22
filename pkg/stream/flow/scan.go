package flow

import (
	"sync"

	"github.com/form3.tech/piper/pkg/stream"
)

// verify scanOperator implements Operator interface
var _ Operator = (*scanOperator)(nil)

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

func scan(zero interface{}, f AggregateFunc) Operator {
	return &scanOperator{
		current: zero,
		f:       f,
	}
}
