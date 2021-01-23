package flow

import (
	"github.com/nicholasham/piper/pkg/piper"
)

type FilterFunc func(value interface{}) bool

// verify filterOperator implements Operator interface
var _ Operator = (*filterOperator)(nil)

type filterOperator struct {
	f FilterFunc
}

func (m *filterOperator) SupportsParallelism() bool {
	return true
}

func (m *filterOperator) Start(actions OperatorActions) {
}

func (m *filterOperator) Apply(element piper.Element, actions OperatorActions) {
	element.WhenValue(func(value interface{}) {
		if m.f(value) {
			actions.PushValue(value)
		}
	}).WhenError(actions.PushError)
}

func (m *filterOperator) End(actions OperatorActions) {
}

func filter(f FilterFunc) Operator {
	return &filterOperator{
		f: f,
	}
}
