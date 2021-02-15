package iterable

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTake(t *testing.T) {

	t.Run("It takes the correct number of elements", func(t *testing.T) {
		assert.Equal(t, Range(0, 4).ToSlice(), Range(0, 50).Take(5).ToSlice())
		assert.Equal(t, Range(0, 9).ToSlice(), Range(0, 50).Take(10).ToSlice())
		assert.Equal(t, Range(1, 8).ToSlice(), Range(1, 50).Take(8).ToSlice())
	})

	t.Run("It returns nil when empty", func(t *testing.T) {
		assert.Equal(t, Empty().ToSlice(),  Empty().Take(20).ToSlice())
	})

}

