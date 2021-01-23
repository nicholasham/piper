package flow

import (
	"github.com/nicholasham/piper/pkg/stream"
)

type MapFunc func(value interface{}) (interface{}, error)

// verify mapOperator implements OperatorLogic interface
var _ OperatorLogic = (*mapOperator)(nil)

type mapOperator struct {
	f MapFunc
}

func (m *mapOperator) SupportsParallelism() bool {
	return true
}

func (m *mapOperator) Start(actions OperatorActions) {
}

func (m *mapOperator) Apply(element stream.Element, actions OperatorActions) {
	value, err := m.f(element.Value())
	if err != nil {
		actions.PushError(err)
		return
	}
	actions.PushValue(value)
}

func (m *mapOperator) End(actions OperatorActions) {
}

func mapOp(f MapFunc) OperatorLogic {
	return &mapOperator{
		f: f,
	}
}
