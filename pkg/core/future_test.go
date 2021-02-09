package core

import (
	"go.uber.org/goleak"
	"testing"
)

func TestName(t *testing.T) {

	defer goleak.VerifyNone(t)

	future := NewFuture(func() Result {
		return Success(10)
	})

	future.Await()
}
