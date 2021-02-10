package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResult(t *testing.T) {

	t.Run("IsOk is true when in success state", func(t *testing.T) {
		 Ok(10).IsOk()
	})

	t.Run("IsErr true when in error state", func(t *testing.T) {
		Err(fmt.Errorf("some error")).IsErr()
	})

	t.Run("Unwrap returns value when in success state ", func(t *testing.T) {
		value, err := Ok(10).Unwrap()
		assert.NotNil(t, value)
		assert.NoError(t, err)
	})

	t.Run("Unwrap returns error when in error state ", func(t *testing.T) {
		value, err := Err(fmt.Errorf("some error")).Unwrap()
		assert.Nil(t, value)
		assert.Error(t, err)
	})

	t.Run("Map returns the mapped ok value in another result", func(t *testing.T) {
		result := Ok(10).Map(func(value Any) Any {
			return value.(int) + 20
		})

		value, err := result.Unwrap()

		assert.Equal(t, 30, value)
		assert.NoError(t, err)
	})

	t.Run("Map returns error result when in error state", func(t *testing.T) {
		result := Err(fmt.Errorf("some error")).Map(func(value Any) Any {
			return value.(int) + 20
		})

		value, err := result.Unwrap()

		assert.Nil(t, value)
		assert.Error(t, err)
	})

	t.Run("Then continues the computation", func(t *testing.T) {
		result := Ok(10).Then(func(value Any) Result {
			return Ok(value.(int) + 20)
		})

		value, err := result.Unwrap()

		assert.Equal(t, 30, value)
		assert.NoError(t, err)
	})

	t.Run("Then does not continue when in error state", func(t *testing.T) {
		result := Err(fmt.Errorf("some error")).Then(func(value Any) Result {
			return Ok(value.(int) + 20)
		})

		value, err := result.Unwrap()

		assert.Nil(t, value)
		assert.Error(t, err)
	})



	t.Run("OrElse handles the error and continues computation", func(t *testing.T) {
		result := Err(fmt.Errorf("some error")).OrElse(func(err error) Result {
			t.Logf("handled: %v", err)
			return Ok(20)
		})

		value, err := result.Unwrap()

		assert.Equal(t, 20, value)
		assert.NoError(t, err)
	})

	t.Run("OrElse does not run error handler when in ok state", func(t *testing.T) {
		result := Ok(50).OrElse(func(err error) Result {
			t.Logf("handled: %v", err)
			return Ok(20)
		})

		value, err := result.Unwrap()

		assert.Equal(t , 50, value)
		assert.NoError(t, err)
	})


		t.Run("Match allows returning result from success or failure ", func(t *testing.T) {

		onSuccess := func(value Any) Any {
			return value.(int) * 10
		}

		onFailure := func(value error) Any {
			return 100000
		}

		goodValue := Ok(10).Match(onSuccess, onFailure)
		badValue := Err(fmt.Errorf("some error")).Match(onSuccess, onFailure)

		assert.Equal(t, 100, goodValue)
		assert.Equal(t, 100000, badValue)

	})

	t.Run("IfOk runs action when in success state", func(t *testing.T) {

		actual := 100

		Ok(20).IfOk(func(value Any) {
			actual = actual + value.(int)
		}).IfErr(func(err error) {
			actual = 0
		})

		assert.Equal(t, 120, actual)
	})

	t.Run("IfErr runs action when in error state", func(t *testing.T) {

		actual := 100

		Err(fmt.Errorf("some error")).IfOk(func(value Any) {
			actual = actual + value.(int)
		}).IfErr(func(err error) {
			actual = 0
		})

		assert.Equal(t, 0, actual)
	})

	}
