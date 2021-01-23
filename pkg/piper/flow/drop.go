package flow

import (
	"github.com/nicholasham/piper/pkg/piper"
	"sync/atomic"

)

// verify dropOperator implements Operator interface
var _ Operator = (*dropOperator)(nil)

type dropOperator struct {
	op     uint64
	number int
}

func (receiver *dropOperator) Start(_ OperatorActions) {
}

func (receiver *dropOperator) Apply(element piper.Element, actions OperatorActions) {
	current := int(atomic.AddUint64(&receiver.op, 1))
	if current > receiver.number {
		element.
			WhenValue(actions.PushValue).
			WhenError(actions.PushError)
	}
	return
}

func (receiver *dropOperator) End(_ OperatorActions) {
}

func (receiver *dropOperator) SupportsParallelism() bool {
	return true
}

func drop(number int) Operator {
	return &dropOperator{number: number}
}
