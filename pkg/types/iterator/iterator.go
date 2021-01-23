package iterator

import "fmt"

type Iterator interface {
	HasNext() bool
	Next() (interface{}, error)
	ToList() []interface{}
}

var EndOfError error = fmt.Errorf("iterator reached end")
var EmptyError error = fmt.Errorf("iterator empty")

func toList(iterator Iterator) []interface{} {
	var values []interface{}
	for iterator.HasNext() {

		value, _ := iterator.Next()
		values = append(values, value)
	}
	return values
}
