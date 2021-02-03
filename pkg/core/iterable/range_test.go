package iterable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRange(t *testing.T) {
	t.Run("Able to iterate over full range", func(t *testing.T) {
		sut := Range(1, 10)

		hasNextCount := 0
		expected := []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		var values []interface{}
		iterator := sut.Iterator()
		for iterator.HasNext() {
			hasNextCount++
			value := iterator.Next()
			values = append(values, value)
		}

		assert.Equal(t, 10, hasNextCount)
		assert.Equal(t, expected, values)
	})
}
