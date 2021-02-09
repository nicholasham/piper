package core

import (
	"go.uber.org/goleak"
	"testing"
)

func TestName(t *testing.T) {
	defer goleak.VerifyNone(t)

	promise := NewPromise()

	go func(p *Promise) {
		p.TrySuccess(10)
	}(promise)

	promise.Future().Await()
}
