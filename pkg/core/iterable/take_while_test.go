package iterable

import (
	"github.com/nicholasham/piper/pkg/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTakeWhile(t *testing.T) {

	numbers := Range(0, 50)

	inRange := func(start int, end int) core.PredicateFunc {
		return func(value core.Any) bool {
			return  value.(int) >= start && value.(int) <= end
		}
	}

	t.Run("It returns the elements from the list as long as the condition is satisfied.", func(t *testing.T) {
		assert.Equal(t, Range(0, 5).ToSlice(), numbers.TakeWhile(inRange(0, 5)).ToSlice())
		assert.Equal(t, Range(0, 10).ToSlice(), numbers.TakeWhile(inRange(0,10)).ToSlice())
	})

	t.Run("It returns nil when not in range", func(t *testing.T) {
		assert.Equal(t, Empty().ToSlice(), numbers.TakeWhile(inRange(51, 52)).ToSlice())
	})



}

