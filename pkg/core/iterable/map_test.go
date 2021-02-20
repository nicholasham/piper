package iterable

import (
	"github.com/nicholasham/piper/pkg/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {

	numbers := Range(1, 10)

	multiplyBy := func(multiple int) core.MapFunc {
		return func(value core.Any) core.Any {
			return value.(int) * multiple
		}
	}

	t.Run("It maps each iterated value", func(t *testing.T) {
		assert.Equal(t, Slice(2, 4, 6, 8, 10, 12, 14, 16, 18, 20).ToSlice(), numbers.Map(multiplyBy(2)).ToSlice())
	})

	t.Run("It panics when calling next with no next value", func(t *testing.T) {
		assert.Panics(t, func() {
			iterator := Range(1,1).Iterator()
			iterator.Next()
			iterator.Next()
		})
	})

}
