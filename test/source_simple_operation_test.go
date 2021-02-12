package test

import (
	"context"
	"github.com/nicholasham/piper/pkg/core/iterable"
	"github.com/nicholasham/piper/pkg/stream"
	"github.com/nicholasham/piper/pkg/stream/sink"
	"github.com/nicholasham/piper/pkg/stream/source"

	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func List(values ...interface{}) []interface{} {
	return values
}

func TestMap(t *testing.T) {
	defer goleak.VerifyNone(t)

	mapping := func(value interface{}) (interface{}, error) {
		return value.(int) * 2, nil
	}

	result := source.
		Slice(1, 2, 3, 4, 5).
		Map(mapping).With(stream.Parallelism(5000)).
		RunWith(context.Background(), sink.Slice())

	values, err := result.Await().Unwrap()

	assert.NoError(t, err)
	assert.EqualValues(t, List(2, 4, 6, 8, 10), values)
}




func TestMapConcat(t *testing.T) {
	defer goleak.VerifyNone(t)

	mapping := func(value interface{}) (iterable.Iterable, error) {
		return iterable.Slice(value, value), nil
	}

	result := source.
		Slice(1, 2, 3, 4, 5).
		MapConcat(mapping).
		To(sink.Slice()).
		Run(context.Background())

	values, err := result.Await().Unwrap()

	assert.NoError(t, err)
	assert.EqualValues(t, []interface{}{1, 1, 2, 2, 3, 3, 4, 4, 5, 5}, values)
}


func TestScan(t *testing.T) {

	addOne := func(acc interface{}, out interface{}) (interface{}, error) {
		return acc.(int) + out.(int), nil
	}

	expected := []interface{} {0, 1, 3, 6, 10, 15}

	result := source.
		Slice(1, 2, 3, 4, 5).
		//Log("Before Scan").
		Scan(0, addOne).
		//Log("After Scan").
		To(sink.Slice()).Run(context.Background())

	values, err := result.Await().Unwrap()

	assert.NoError(t, err)
	assert.ElementsMatch(t, expected, values)

}


func TestFold(t *testing.T) {

	sum := func(acc interface{}, out interface{}) (interface{}, error) {
		return acc.(int) + 1, nil
	}

	expected := 100000

	result := source.
		Range(1, expected).
		Fold(0, sum).
		To(sink.Head()).
		Run(context.Background())

	value, err := result.Await().Unwrap()

	assert.NoError(t, err)
	assert.EqualValues(t, expected, value)

}

func TestFilter(t *testing.T) {
	defer goleak.VerifyNone(t)

	evenNumbersOnly := func(value interface{}) bool {
		return value.(int)%2 == 0
	}

	expected := []interface{}{2, 4, 6, 8, 10}

	result := source.
		Slice(1, 2, 3, 4, 5, 6, 7, 8, 9, 10).
		Filter(evenNumbersOnly).
		To(sink.Slice()).
		Run(context.Background())

	values, err := result.Await().Unwrap()

	assert.NoError(t, err)
	assert.ElementsMatch(t, expected, values)

}

func TestTake(t *testing.T) {
	defer goleak.VerifyNone(t)

	expected := [] interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := source.
		Range(1, 100000000).
		Take(10).
		To(sink.Slice()).
		Run(context.Background())

	values, err := result.Await().Unwrap()

	assert.NoError(t, err)
	assert.ElementsMatch(t, expected, values)
}
