package iterable

import "github.com/nicholasham/piper/pkg/core"

// Provides a set of options to create iterable values
type Iterable interface {
	// Creates a new tail over all elements contained in this iterable object.
	Iterator() Iterator
	// Applies a function f to all items in the iterable
	ForEach(f func(item interface{}))

	// Appends all items in the iterable to a slice
	ToSlice() []interface{}

	// Creates an iterable that iterates items while teh p is satisfied
	TakeWhile(f  core.PredicateFunc ) Iterable
}

// verify iterable implements Iterable interface
var _ Iterable = (*iterable)(nil)

type iterable struct {
	newIterator func() Iterator
}

func (i *iterable) ToSlice() []interface{} {
	var items []interface{}
	i.ForEach(func(item interface{}) {
		items = append(items, item)
	})
	return items
}

func NewIterable(f func() Iterator) Iterable {
	return &iterable{newIterator: f}
}

func (i *iterable) Iterator() Iterator {
	return i.newIterator()
}

func (i *iterable) ForEach(f func(item interface{})) {
	iterator := i.newIterator()
	for iterator.HasNext() {
		f(iterator.Next())
	}
}

func  (i *iterable)  TakeWhile(f  core.PredicateFunc ) Iterable {
	return takeWhile(f, i.Iterator())
}
