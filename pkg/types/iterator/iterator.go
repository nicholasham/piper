package iterator

import "fmt"

type T interface{}

type Iterator interface {
	HasNext() bool
	Next() (T, error)
	ToList() []T
}

var EndOfError error = fmt.Errorf("iterator reached end")
var EmptyError error = fmt.Errorf("iterator empty")


func toList(iterator Iterator) []T {
	var values []T
	for iterator.HasNext() {

		value, _ := iterator.Next()
		values = append(values, value)
	}
	return values
}

