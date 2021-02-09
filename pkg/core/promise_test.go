package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"testing"
)

func TestPromise(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("can only write success once ", func(t *testing.T) {
		promise := NewPromise()
		assert.True(t, promise.TrySuccess(10))
		assert.False(t, promise.TrySuccess(20))
		assert.False(t, promise.TryFailure(fmt.Errorf("some error")))
		assert.Equal(t, Success(10), <-promise.resultChan)
	})

	t.Run("can only write failure once ", func(t *testing.T) {
		promise := NewPromise()
		expectedError := fmt.Errorf("some error")
		ignoredError := fmt.Errorf("some error")

		assert.True(t, promise.TryFailure(expectedError))
		assert.False(t, promise.TryFailure(ignoredError))
		assert.False(t, promise.TrySuccess(20))
		assert.Equal(t, Failure(expectedError), <-promise.resultChan)
	})

}
