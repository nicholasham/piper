package iterable

// verify emptyIterator implements Iterator interface
var _ Iterator = (*emptyIterator)(nil)


type emptyIterator struct {
}

func (e *emptyIterator) HasNext() bool {
	return false
}

func (e *emptyIterator) Next() interface{} {
	panic("next was called on empty iterator")
}

func Empty() Iterable {
	return NewIterable(func() Iterator {
		return &emptyIterator{}
	})
}
