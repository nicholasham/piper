package iterable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {

	t.Run("HasNext always yields false", func(t *testing.T) {
		sut := Empty()
		assert.False(t, sut.Iterator().HasNext())
	})

	t.Run("Next always panics", func(t *testing.T) {
		sut := Empty()
		assert.Panics(t, func() {
			sut.Iterator().Next()
		})
	})

}
