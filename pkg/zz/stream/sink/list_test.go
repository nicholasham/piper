package sink

import (
	"context"
	"testing"

	"github.com/nicholasham/piper/pkg/types/iterator"
	"github.com/nicholasham/piper/pkg/zz/stream/source"

	"go.uber.org/goleak"

	"github.com/stretchr/testify/assert"
)

func TestSliceSink(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("Must return a slice from a source", func(t *testing.T) {
		var input = iterator.Range(1, 1000, 1).ToList()
		future := source.List(input).RunWith(context.Background(), List())
		value, error := future.Await()
		assert.NoError(t, error)
		assert.Equal(t, input, value)
	})

	t.Run("Must return an empty sequence from an empty source", func(t *testing.T) {
		var input = iterator.Empty().ToList()
		future := source.List(input).RunWith(context.Background(), List())
		value, error := future.Await()
		assert.NoError(t, error)
		assert.Equal(t, input, value)
	})

	t.Run("Must fail on cancellation", func(t *testing.T) {
		var input = iterator.Range(1, 1000, 1).ToList()
		ctx, cancel := context.WithCancel(context.Background())
		future := source.List(input).RunWith(ctx, List())
		cancel()
		_, err := future.Await()
		assert.Error(t, err)
		assert.EqualError(t, err, context.Canceled.Error())
	})
}
