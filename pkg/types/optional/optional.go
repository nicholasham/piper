package optional

import (
	"fmt"
	"github.com/nicholasham/piper/pkg/types/iterator"
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
	if o.IsEmpty() {
		return nil, EmptyError
	}
	return o.value, nil
}

// Returns true if the option is None, false otherwise.
func (o Option) IsEmpty() bool {
	return !o.hasValue
}

// Returns true if this option is nonempty and the predicate p returns true when applied to this Option's value. Otherwise, returns false.
func (o Option) Exists(f PredicateFunc) bool {
	if !o.IsDefined() {
		return false
	}
	return f(o.value)
}

// Returns true if the option is an instance of Some, false otherwise.
func (o Option) IsDefined() bool {
	return o.hasValue
}

// Returns the option's value if the option is nonempty, otherwise return the result of evaluating default.
func (o Option) GetOrElse(defaultValue interface{}) interface{} {
	if o.IsEmpty() {
		return defaultValue
	}
	return o.value
}

// Returns this Option if it is nonempty and applying the predicate p to this Option's value returns true. Otherwise, return None.
func (o Option) Filter(f PredicateFunc) Option {
	if o.IsDefined() && o.Exists(f) {
		return o
	}
	return None()
}

// Returns the result of applying f to this Option's value if this Option is nonempty. Returns None if this Option is empty.
func (o Option) FlatMap(f FlatMapFunc) Option {
	if o.IsEmpty() {
		return o
	}
	return f(o.value)
}

// Returns a Some containing the result of applying f to this Option's value if this Option is nonempty. Otherwise return None.
func (o Option) Map(f MapFunc) Option {
	if o.IsEmpty() {
		return o
	}
	return Some(f(o.value))
}

func (o Option) Iterator() iterator.Iterator {
	if o.IsEmpty() {
		return iterator.Empty()
	}
	return iterator.Single(o.value)
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
