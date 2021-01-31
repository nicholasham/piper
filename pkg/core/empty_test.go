package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {

	t.Run("HasNext is always false", func(t *testing.T) {
		sut := Empty()
		assert.False(t, sut.Iterator().HasNext())
	})

	t.Run("Next returns empty error", func(t *testing.T) {
		sut := Empty()
		value := sut.Iterator().Next()
		assert.Nil(t, value)
	})

}
