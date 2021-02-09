package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptional(t *testing.T) {

	var add = func(number int) func(value Any) Any {
		return func(value Any) Any {
			return value.(int) + number
		}
	}

	var isEqualTo = func(expected Any) func(actual Any) bool {
		return func(actual Any) bool {
			return expected == actual
		}
	}

	t.Run("is structurally equal", func(t *testing.T) {
		assert.Equal(t, Some(1), Some(1))
		assert.Equal(t, None(), None())
		assert.NotEqual(t, Some(1), Some(2))
	})

	t.Run("IsSome is true when in some state", func(t *testing.T) {
		assert.True(t, Some(1).IsSome())
		assert.False(t, None().IsSome())
	})

	t.Run("Await returns value when has some value", func(t *testing.T) {
		expected := 100
		actual, err := Some(100).Get()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Await returns error when none", func(t *testing.T) {

		actual, err := None().Get()

		assert.Nil(t, actual)
		assert.Error(t, err)
		assert.Equal(t, EmptyError, err)
	})

	t.Run("GetOrElse returns value when has some value", func(t *testing.T) {
		assert.Equal(t, 100, Some(100).GetOrElse(200))
	})

	t.Run("GetOrElse returns default value when nothing", func(t *testing.T) {
		assert.Equal(t, 200, None().GetOrElse(200))
	})

	t.Run("IsNone is only true for none", func(t *testing.T) {
		assert.True(t, None().IsNone())
		assert.False(t, Some(1).IsNone())
	})

	t.Run("Exists returns true if the predicate matches the value ", func(t *testing.T) {
		assert.True(t, Some(1).Exists(func(value Any) bool {
			return value.(int) == 1
		}))
	})

	t.Run("Exists returns false if the predicate does not match the value ", func(t *testing.T) {
		assert.False(t, Some(1).Exists(func(value Any) bool {
			return value.(int) == 2
		}))
	})

	t.Run("Exists returns false when none ", func(t *testing.T) {
		assert.False(t, None().Exists(func(value Any) bool {
			return true
		}))
	})

	t.Run("Map returns wrapped mapped value for something", func(t *testing.T) {
		expected := Some(2)
		actual := Some(1).Map(add(1))
		assert.Equal(t, expected, actual)
	})

	t.Run("Map returns none when none", func(t *testing.T) {
		expected := None()
		actual := expected.Map(add(1))
		assert.Equal(t, expected, actual)
	})

	t.Run("FlatMap returns mapped something", func(t *testing.T) {
		expected := Some(2)
		actual := Some(1).FlatMap(func(value Any) Optional {
			return Some(2)
		})
		assert.Equal(t, expected, actual)
	})

	t.Run("FlatMap returns none when none", func(t *testing.T) {
		expected := None()
		actual := expected.FlatMap(func(value Any) Optional {
			return Some(1)
		})
		assert.Equal(t, expected, actual)
	})

	t.Run("Filter returns none when none", func(t *testing.T) {
		expected := None()
		actual := expected.Filter(isEqualTo(100))
		assert.Equal(t, expected, actual)
	})

	t.Run("Filter returns some when has value that matches predicate", func(t *testing.T) {
		expected := Some(100)
		actual := expected.Filter(isEqualTo(100))
		assert.Equal(t, expected, actual)
	})

	t.Run("Filter returns none when has value that does not match predicate", func(t *testing.T) {
		expected := None()
		actual := Some(200).Filter(isEqualTo(100))
		assert.Equal(t, expected, actual)
	})

}
