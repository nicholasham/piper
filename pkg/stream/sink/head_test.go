package sink

import (
	"context"
	"fmt"
	"go.uber.org/goleak"
	"testing"

	"github.com/nicholasham/piper/pkg/stream/source"

	"github.com/stretchr/testify/assert"
)

func TestHead(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Must yield first value", func(t *testing.T) {
		future := source.Range(1, 100).RunWith(context.Background(), Head())
		value, error := future.Await().Unwrap()
		assert.NoError(t, error)
		assert.Equal(t, 1, value)
	})

	t.Run("Must yield error when empty source", func(t *testing.T) {
		future := source.Empty().RunWith(context.Background(), Head())
		value, error := future.Await().Unwrap()
		assert.Error(t, error)
		assert.Nil(t, value)
		assert.Equal(t, HeadOfEmptyStream, error)
	})

	t.Run("Must yield error when error in source", func(t *testing.T) {
		expectedErr := fmt.Errorf("some error occured")
		future := source.Failed(expectedErr).RunWith(context.Background(), Head())
		value, error := future.Await().Unwrap()
		assert.Error(t, error)
		assert.Nil(t, value)
		assert.Equal(t, expectedErr, error)
	})

}
