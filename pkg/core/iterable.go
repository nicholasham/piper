package core

// Provides a set of options to create iterable values
type Iterable interface {
	// Creates a new iterator over all elements contained in this iterable object.
	Iterator() Iterator
	// Applies a function f to all items in the iterable
	ForEach(f func(item interface{}))
}

// verify iterable implements Iterable interface
var _ Iterable = (*iterable)(nil)

type iterable struct {
	newIterator func() Iterator
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


