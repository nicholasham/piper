package iterable

import (
	"fmt"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTakeWhile(t *testing.T) {

	numbers := Range(0, 50)
	zeroTo9 := numbers.TakeWhile(func(value core.Any) bool {
		return value.(int) < 10
	})

	zeroTo9.ForEach(func(item interface{}) {
		fmt.Println(item)
	})


	zeroTo1 := numbers.TakeWhile(func(value core.Any) bool {
		return value.(int) < 2
	}).Iterator()

	assert.True(t, zeroTo1.HasNext())
	assert.Equal(t, 0, zeroTo1.Next())
	assert.True(t, zeroTo1.HasNext())
	assert.Equal(t, 1, zeroTo1.Next())




	t.Run("Able to iterate over all items in slice", func(t *testing.T) {

		sut := Slice(1,2,3,4,5,6,7,8,9, 10).TakeWhile(func(value core.Any) bool {
			result := value.(int) <=5
			return result
		})

		hasNextCount := 0
		expected := []interface{}{1, 2, 3, 4, 5}
		var values []interface{}
		iterator := sut.Iterator()
		for iterator.HasNext() {
			hasNextCount++
			value := iterator.Next()
			values = append(values, value)
		}

		assert.Equal(t, 5, hasNextCount)
		assert.Equal(t, expected, values)
	})

}