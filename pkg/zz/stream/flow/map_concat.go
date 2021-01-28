package flow

import (
	"github.com/nicholasham/piper/pkg/zz/stream"
)

type MapConcatFunc func(value interface{}) ([]interface{}, error)

// verify mapConcatOperator implements OperatorLogic interface
var _ OperatorLogic = (*mapConcatOperator)(nil)

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

func mapConcat(f MapConcatFunc) OperatorLogic {
	return &mapConcatOperator{
		f: f,
	}
}
