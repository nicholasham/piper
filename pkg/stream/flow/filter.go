package flow

import (
	"github.com/nicholasham/piper/pkg/stream"
)

type FilterFunc func(value interface{}) bool

// verify filterOperator implements OperatorLogic interface
var _ OperatorLogic = (*filterOperator)(nil)

type filterOperator struct {
	f FilterFunc
}

func (m *filterOperator) SupportsParallelism() bool {
	return true
}

func (m *filterOperator) Start(actions OperatorActions) {
}

func (m *filterOperator) Apply(element stream.Element, actions OperatorActions) {
	element.WhenValue(func(value interface{}) {
		if m.f(value) {
			actions.PushValue(value)
		}
	}).WhenError(actions.PushError)
}

func (m *filterOperator) End(actions OperatorActions) {
}

func filter(f FilterFunc) OperatorLogic {
	return &filterOperator{
		f: f,
	}
}
