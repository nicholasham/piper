package flow

import (
	"sync/atomic"

	"github.com/form3.tech/piper/pkg/stream"
)

// verify takeOperator implements Operator interface
var _ Operator = (*takeOperator)(nil)

type takeOperator struct {
	op     uint64
	number int
}

func (receiver *takeOperator) Start(actions OperatorActions) {
}

func (receiver *takeOperator) Apply(element stream.Element, actions OperatorActions) {
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

func take(number int) Operator {
	return &takeOperator{number: number}
}
