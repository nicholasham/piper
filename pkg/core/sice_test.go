package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	t.Run("Able to iterate over all items in slice", func(t *testing.T) {
		sut := Slice(1, 2, 3)

		hasNextCount := 0
		expected := []interface{}{1, 2, 3}
		var values []interface{}
		iterator := sut.Iterator()
		for iterator.HasNext() {
			hasNextCount++
			value := iterator.Next()
			values = append(values, value)
		}

		assert.Equal(t, 3, hasNextCount)
		assert.Equal(t, expected, values)
	})
}
