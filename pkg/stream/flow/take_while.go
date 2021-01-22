package flow

import (
	"github.com/form3.tech/piper/pkg/stream"
)

// verify takeWhileOperator implements Operator interface
var _ Operator = (*takeWhileOperator)(nil)

type takeWhileOperator struct {
	f FilterFunc
}

func (receiver *takeWhileOperator) Start(actions OperatorActions) {
}

func (receiver *takeWhileOperator) Apply(element stream.Element, actions OperatorActions) {

	element.WhenError(actions.PushError)
	element.WhenValue(func(value interface{}) {
		if !receiver.f(value) {
			actions.CompleteStage()
		}
		actions.PushValue(value)
	})
}

func (receiver *takeWhileOperator) End(actions OperatorActions) {
}

func (receiver *takeWhileOperator) SupportsParallelism() bool {
	return true
}

func takeWhile(f FilterFunc) Operator {
	return &takeWhileOperator{f: f}
}
