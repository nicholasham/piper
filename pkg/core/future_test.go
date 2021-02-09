package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"testing"
)

func TestAwait(t *testing.T) {

	defer goleak.VerifyNone(t)

	t.Run("Test success result", func(t *testing.T) {
		future := NewFuture(func() Result {
			return Success(10)
		})

		result := future.Await()
		assert.Equal(t, Success(10), result)
	})

	t.Run("Test failure result", func(t *testing.T) {
		future := NewFuture(func() Result {
			return Failure(fmt.Errorf("some error"))
		})

		result := future.Await()
		assert.Equal(t, Failure(fmt.Errorf("some error")), result)
	})


}
