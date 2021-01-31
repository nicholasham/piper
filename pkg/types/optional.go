package types

import (
	"fmt"
)

// Is a container for zero or one element of a given type
type Optional struct {
	value    T
	hasValue bool
}

type FlatMapFunc func(value T) Optional
type MapFunc func(value T) T
type PredicateFunc func(value T) bool

type MapSome func(value T) R
type MapNone func() R

type T interface{}

var EmptyError = fmt.Errorf("is empty")

// Returns the option's value.
func (o Optional) Get() (T, error) {
	if o.IsNone() {
		return nil, EmptyError
	}
	return o.value, nil
}

// Returns true if the option is in a None state, false otherwise.
func (o Optional) IsNone() bool {
	return !o.hasValue
}

// Returns true if this option is nonempty and the predicate p returns true when applied to this Optional's value. Otherwise, returns false.
func (o Optional) Exists(f PredicateFunc) bool {
	if !o.IsSome() {
		return false
	}
	return f(o.value)
}

// Returns true if the option is a Some state, false otherwise.
func (o Optional) IsSome() bool {
	return o.hasValue
}

func (o Optional) IfSome(f func(value interface{})) Optional {
	if o.IsSome() {
		f(o.value)
	}
	return o
}

func (o Optional) IfNone(f func()) Optional {
	if o.IsNone() {
		f()
	}
	return o
}

// Returns the option's value if the option is nonempty, otherwise return the result of evaluating default.
func (o Optional) GetOrElse(defaultValue interface{}) interface{} {
	if o.IsNone() {
		return defaultValue
	}
	return o.value
}

// Returns this Optional if it is nonempty and applying the predicate p to this Optional's value returns true. Otherwise, return None.
func (o Optional) Filter(f PredicateFunc) Optional {
	if o.IsSome() && o.Exists(f) {
		return o
	}
	return None()
}

// Returns the result of applying f to this Optional's value if this Optional is nonempty. Returns None if this Optional is empty.
func (o Optional) FlatMap(f FlatMapFunc) Optional {
	if o.IsNone() {
		return o
	}
	return f(o.value)
}

// Returns a Some containing the result of applying f to this Optional's value if this Optional is nonempty. Otherwise return None.
func (o Optional) Map(f MapFunc) Optional {
	if o.IsNone() {
		return o
	}
	return Some(f(o.value))
}

func (o Optional) Match(some MapSome, none MapNone) R {
	if o.IsNone() {
		return none()
	}
	return some(o.value)
}

func None() Optional {
	return Optional{
		value:    nil,
		hasValue: false,
	}
}

func Some(value interface{}) Optional {
	return Optional{
		value:    value,
		hasValue: true,
	}
}
