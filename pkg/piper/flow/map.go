package flow

import (
	"github.com/nicholasham/piper/pkg/piper"
)

type MapFunc func(value interface{}) (interface{}, error)

// verify mapOperator implements Operator interface
var _ Operator = (*mapOperator)(nil)

type mapOperator struct {
	f MapFunc
}

func (m *mapOperator) SupportsParallelism() bool {
	return true
}

func (m *mapOperator) Start(actions OperatorActions) {
}

func (m *mapOperator) Apply(element piper.Element, actions OperatorActions) {
	value, err := m.f(element.Value())
	if err != nil {
		actions.PushError(err)
		return
	}
	actions.PushValue(value)
}

func (m *mapOperator) End(actions OperatorActions) {
}

func mapOp(f MapFunc) Operator {
	return &mapOperator{
		f: f,
	}
}
