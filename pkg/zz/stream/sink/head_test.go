package sink

import (
	"context"
	"fmt"
	"testing"

	"github.com/nicholasham/piper/pkg/zz/stream/source"

	"go.uber.org/goleak"

	"github.com/stretchr/testify/assert"
)

func TestHead(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Must yield first value", func(t *testing.T) {
		future := source.Range(1, 100000000, 1).RunWith(context.Background(), Head())
		value, error := future.Await()
		assert.NoError(t, error)
		assert.Equal(t, 1, value)
	})

	t.Run("Must yield first error", func(t *testing.T) {
		cause := fmt.Errorf("must fail and return this error")
		future := source.Failed(cause).RunWith(context.Background(), Head())
		value, error := future.Await()
		assert.Error(t, error)
		assert.Nil(t, value)
		assert.Equal(t, cause, error)
	})

}
