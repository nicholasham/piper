package iterable

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDrop(t *testing.T) {

	t.Run("It drops n items before returning the correct items", func(t *testing.T) {
		assert.Equal(t, Range(5, 50).ToSlice(), Range(0, 50).Drop(5).ToSlice())
		assert.Equal(t, Range(10, 50).ToSlice(), Range(0, 50).Drop(10).ToSlice())
		assert.Equal(t, Range(9, 50).ToSlice(), Range(1, 50).Drop(8).ToSlice())
	})

	t.Run("It returns nil when empty", func(t *testing.T) {
		assert.Equal(t, Empty().ToSlice(),  Empty().Drop(20).ToSlice())
	})

}

