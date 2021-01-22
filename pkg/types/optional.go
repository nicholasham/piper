package types

import (
	"fmt"
)

// Is a container for zero or one element of a given type
type Option struct {
	value    T
	hasValue bool
}

type FlatMapFunc func(value T) Option
type MapFunc func(value T) T
type PredicateFunc func(value T) bool

type T interface{}

var EmptyError = fmt.Errorf("is empty")

// Returns the option's value.
func (o Option) Get() (T, error) {
	if o.IsNone() {
		return nil, EmptyError
	}
	return o.value, nil
}

// Returns true if the option is in a None state, false otherwise.
func (o Option) IsNone() bool {
	return !o.hasValue
}

// Returns true if this option is nonempty and the predicate p returns true when applied to this Option's value. Otherwise, returns false.
func (o Option) Exists(f PredicateFunc) bool {
	if !o.IsSome() {
		return false
	}
	return f(o.value)
}

// Returns true if the option is a Some state, false otherwise.
func (o Option) IsSome() bool {
	return o.hasValue
}

// Returns the option's value if the option is nonempty, otherwise return the result of evaluating default.
func (o Option) GetOrElse(defaultValue interface{}) interface{} {
	if o.IsNone() {
		return defaultValue
	}
	return o.value
}

// Returns this Option if it is nonempty and applying the predicate p to this Option's value returns true. Otherwise, return None.
func (o Option) Filter(f PredicateFunc) Option {
	if o.IsSome() && o.Exists(f) {
		return o
	}
	return None()
}

// Returns the result of applying f to this Option's value if this Option is nonempty. Returns None if this Option is empty.
func (o Option) FlatMap(f FlatMapFunc) Option {
	if o.IsNone() {
		return o
	}
	return f(o.value)
}

// Returns a Some containing the result of applying f to this Option's value if this Option is nonempty. Otherwise return None.
func (o Option) Map(f MapFunc) Option {
	if o.IsNone() {
		return o
	}
	return Some(f(o.value))
}

func None() Option {
	return Option{
		value:    nil,
		hasValue: false,
	}
}

func Some(value interface{}) Option {
	return Option{
		value:    value,
		hasValue: true,
	}
}
