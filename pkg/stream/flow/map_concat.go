package flow

import (
	"github.com/form3.tech/piper/pkg/stream"
)

type MapConcatFunc func(value interface{}) ([]interface{}, error)

// verify mapConcatOperator implements Operator interface
var _ Operator = (*mapConcatOperator)(nil)

type mapConcatOperator struct {
	f MapConcatFunc
}

func (m *mapConcatOperator) SupportsParallelism() bool {
	return true
}

func (m *mapConcatOperator) Start(_ OperatorActions) {
}

func (m *mapConcatOperator) Apply(element stream.Element, actions OperatorActions) {
	element.WhenValue(func(value interface{}) {
		values, err := m.f(value)
		if err != nil {
			actions.PushError(err)
			return
		}

		for _, value := range values {
			actions.PushValue(value)
		}

	})
}

func (m *mapConcatOperator) End(actions OperatorActions) {

}

func mapConcat(f MapConcatFunc) Operator {
	return &mapConcatOperator{
		f: f,
	}
}
