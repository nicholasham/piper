package stream

import "fmt"

type MapValue func(value interface{}) (interface{}, error)
type MapError func(err error) error

type ValueAction func(value interface{})
type ErrorAction func(err error)

type Element struct {
	value interface{}
	error error
}

func ToElement(f func() (interface{}, error)) Element {
	value, err := f()
	if err != nil {
		return Error(err)
	} else {
		return Value(value)
	}
}

var NotUsed Element = Element{}

func Value(value interface{}) Element {
	return Element{value: value}
}

func Error(err error) Element {
	return Element{error: err}
}

func Errorf(format string, a ...interface{}) Element {
	return Error(fmt.Errorf(format, a...))
}

func IsError(element Element) bool {
	return element.IsError()
}

func IsValue(element Element) bool {
	return !element.IsError()
}

func (e Element) Error() error {
	return e.error
}

func (e Element) Value() interface{} {
	return e.value
}

func (e Element) IsError() bool {
	return e.error != nil
}

func (e Element) IsValue() bool {
	return e.value != nil
}

func Values(values ...interface{}) []Element {
	var elements []Element
	for _, value := range values {
		elements = append(elements, Value(value))
	}
	return elements
}

// Maps the result inputStage the element if it has a result
func (e Element) MapValue(f MapValue) Element {
	if e.IsError() {
		return e
	}
	value, err := f(e.value)
	if err != nil {
		return Error(err)
	}
	return Value(value)
}

// Maps the err inputStage the element if has error
func (e Element) MapError(f MapError) Element {
	if e.IsError() {
		return Error(f(e.error))
	}
	return e
}

func (e Element) Apply(valueAction ValueAction, errorAction ErrorAction) {
	e.WhenValue(valueAction).WhenError(errorAction)
}

func (e Element) WhenError(action ErrorAction) Element {
	if e.IsError() {
		action(e.error)
	}
	return e
}

func (e Element) WhenValue(action ValueAction) Element {
	if !e.IsError() {
		action(e.value)
	}
	return e
}
