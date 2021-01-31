package flow

import (
	"github.com/nicholasham/piper/pkg/zz/stream"
)

// verify takeWhileOperator implements OperatorLogic interface
var _ OperatorLogic = (*takeWhileOperator)(nil)

type takeWhileOperator struct {
	f FilterFunc
}

func (receiver *takeWhileOperator) Start(actions OperatorActions) {
}

func (receiver *takeWhileOperator) Apply(element stream.Element, actions OperatorActions) {

	element.WhenError(actions.SendError)
	element.WhenValue(func(value interface{}) {
		if !receiver.f(value) {
			actions.CompleteStage()
		}
		actions.SendValue(value)
	})
}

func (receiver *takeWhileOperator) End(actions OperatorActions) {
}

func (receiver *takeWhileOperator) SupportsParallelism() bool {
	return true
}

func takeWhile(f FilterFunc) OperatorLogic {
	return &takeWhileOperator{f: f}
}
