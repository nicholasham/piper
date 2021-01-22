package iterator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmpty(t *testing.T) {

	t.Run("HasNext is always false", func(t *testing.T) {
		sut := Empty()
		assert.False(t,  sut.HasNext())
	})

	t.Run("Next returns empty error", func(t *testing.T) {
		sut := Empty()
		value, err := sut.Next()
		assert.Error (t, err)
		assert.Equal(t, EmptyError,  err)
		assert.Nil(t, value)
	})

	t.Run("ToList returns empty slice", func(t *testing.T) {
		sut := Empty()
		assert.Empty(t, sut.ToList())
	})
}
