package flow

import (
	"sync/atomic"

	"github.com/nicholasham/piper/pkg/zz/stream"
)

// verify dropOperator implements OperatorLogic interface
var _ OperatorLogic = (*dropOperator)(nil)

type dropOperator struct {
	op     uint64
	number int
}

func (receiver *dropOperator) Start(_ OperatorActions) {
}

func (receiver *dropOperator) Apply(element stream.Element, actions OperatorActions) {
	current := int(atomic.AddUint64(&receiver.op, 1))
	if current > receiver.number {
		element.
			WhenValue(actions.SendValue).
			WhenError(actions.SendError)
	}
	return
}

func (receiver *dropOperator) End(_ OperatorActions) {
}

func (receiver *dropOperator) SupportsParallelism() bool {
	return true
}

func drop(number int) OperatorLogic {
	return &dropOperator{number: number}
}
