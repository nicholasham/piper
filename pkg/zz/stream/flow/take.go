package flow

import (
	"sync/atomic"

	"github.com/nicholasham/piper/pkg/streamold"
)

// verify takeOperator implements OperatorLogic interface
var _ OperatorLogic = (*takeOperator)(nil)

type takeOperator struct {
	op     uint64
	number int
}

func (receiver *takeOperator) Start(actions OperatorActions) {
}

func (receiver *takeOperator) Apply(element streamold.Element, actions OperatorActions) {
	current := int(atomic.AddUint64(&receiver.op, 1))
	if current <= receiver.number {
		element.
			WhenError(actions.PushError).
			WhenValue(actions.PushValue)
		if current == receiver.number {
			actions.CompleteStage()
		}
	}
}

func (receiver *takeOperator) End(actions OperatorActions) {
}

func (receiver *takeOperator) SupportsParallelism() bool {
	return true
}

func take(number int) OperatorLogic {
	return &takeOperator{number: number}
}
