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

func TestMap(t *testing.T) {
	defer goleak.VerifyNone(t)

	mapping := func(value interface{}) (interface{}, error) {
		return value.(int) * 2, nil
	}

	result := source.
		List(1, 2, 3, 4, 5).
		Map(mapping).With(stream.Parallelism(5000)).
		RunWith(context.Background(), sink.Head())

	values, err := result.Await().Unwrap()

	assert.NoError(t, err)
	assert.EqualValues(t, iterable.Slice(2, 4, 6, 8, 10), values)
}

//
//func TestMapWithWorkerPool(t *testing.T) {
//	defer goleak.VerifyNone(t)
//
//	mapping := func(value interface{}) (interface{}, error) {
//		return value.(int) * 2, nil
//	}
//
//	size := 10000
//	expected := iterator.Range(1, size, 2).ToList()
//
//	result := source.
//		Range(1, size, 1, attribute.OutputBuffer(100)).
//		Via(flow.Map(mapping, attribute.Parallelism(100))).
//		To(sink.List()).
//		Run(context.Background())
//
//	values, err := result.Await()
//
//	assert.NoError(t, err)
//	assert.ElementsMatch(t, expected, values)
//}
//
//func TestMapConcat(t *testing.T) {
//	defer goleak.VerifyNone(t)
//
//	mapping := func(value interface{}) ([]interface{}, error) {
//		return of.Values(value, value), nil
//	}
//
//	result := source.
//		List(of.Integers(1, 2, 3, 4, 5)).
//		Via(flow.MapConcat(mapping)).
//		To(sink.List()).
//		Run(context.Background())
//
//	values, err := result.Await()
//
//	assert.NoError(t, err)
//	assert.EqualValues(t, of.Integers(1, 1, 2, 2, 3, 3, 4, 4, 5, 5), values)
//}
//
//func TestScan(t *testing.T) {
//
//	addOne := func(acc interface{}, out interface{}) (interface{}, error) {
//		return acc.(int) + out.(int), nil
//	}
//
//	expected := of.Integers(0, 1, 3, 6, 10, 15)
//
//	result := source.
//		List(of.Integers(1, 2, 3, 4, 5)).
//		//Log("Before Scan").
//		Via(flow.Scan(0, addOne)).
//		//Log("After Scan").
//		To(sink.List()).Run(context.Background())
//
//	values, err := result.Await()
//
//	assert.NoError(t, err)
//	assert.ElementsMatch(t, expected, values)
//
//}
//
//func TestFold(t *testing.T) {
//
//	sum := func(acc interface{}, out interface{}) (interface{}, error) {
//		return acc.(int) + 1, nil
//	}
//
//	expected := 100
//
//	result := source.
//		Range(1, expected, 1).
//		Via(flow.Log("Before Fold")).
//		Via(flow.Fold(0, sum)).
//		Via(flow.Log("After Fold")).
//		To(sink.Head()).
//		Run(context.Background())
//
//	value, err := result.Await()
//
//	assert.NoError(t, err)
//	assert.EqualValues(t, expected, value)
//
//}
//
//func TestFilter(t *testing.T) {
//	defer goleak.VerifyNone(t)
//
//	evenNumbersOnly := func(value interface{}) bool {
//		return value.(int)%2 == 0
//	}
//
//	expected := of.Integers(2, 4, 6, 8, 10)
//
//	result := source.
//		List(of.Integers(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)).
//		Via(flow.Filter(evenNumbersOnly)).
//		To(sink.List()).
//		Run(context.Background())
//
//	values, err := result.Await()
//
//	assert.NoError(t, err)
//	assert.ElementsMatch(t, expected, values)
//
//}
//
//func TestTake(t *testing.T) {
//	defer goleak.VerifyNone(t)
//
//	expected := of.Integers(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
//
//	result := source.
//		Range(1, 100000000, 1, attribute.OutputBuffer(100)).
//		Via(flow.Take(10)).
//		To(sink.List()).
//		Run(context.Background())
//
//	values, err := result.Await()
//
//	assert.NoError(t, err)
//	assert.ElementsMatch(t, expected, values)
//}
