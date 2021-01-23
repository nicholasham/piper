package flow

import (
	"github.com/nicholasham/piper/pkg/piper"
	"sync"

)

// verify foldOperator implements Operator interface
var _ Operator = (*foldOperator)(nil)

type AggregateFunc func(acc interface{}, out interface{}) (interface{}, error)

type foldOperator struct {
	current interface{}
	f       AggregateFunc
	sync.RWMutex
}

func (receiver *foldOperator) Start(actions OperatorActions) {
}

func (receiver *foldOperator) Apply(element piper.Element, actions OperatorActions) {
	receiver.Lock()
	defer receiver.Unlock()
	out, err := receiver.f(receiver.current, element.Value())
	if err != nil {
		actions.FailStage(err)
	}
	receiver.current = out
}

func (receiver *foldOperator) End(actions OperatorActions) {
	actions.PushValue(receiver.current)
}

func (receiver *foldOperator) SupportsParallelism() bool {
	return true
}

func fold(zero interface{}, f AggregateFunc) Operator {
	return &foldOperator{
		current: zero,
		f:       f,
	}
}
