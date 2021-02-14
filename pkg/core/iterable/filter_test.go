package iterable

import (
	"github.com/nicholasham/piper/pkg/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {

	numbers := Range(0, 10)

	inMultiplesOf := func(multiple int) core.PredicateFunc {
		return func(value core.Any) bool {
			return  value.(int) % multiple == 0
		}
	}

	t.Run("It returns the elements from the list as long as the condition is satisfied.", func(t *testing.T) {
		assert.Equal(t, Slice(2, 4,6,8,10).ToSlice(), numbers.Filter(inMultiplesOf(2)).ToSlice())
	})

	t.Run("It returns nil when nothing matches filter", func(t *testing.T) {
		assert.Equal(t, Empty().ToSlice(), numbers.Filter(inMultiplesOf(100)).ToSlice())
	})



}

