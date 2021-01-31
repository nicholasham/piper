package core

// verify emptyIterator implements Iterator interface
var _ Iterator = (*emptyIterator)(nil)

type emptyIterator struct {
}

func (e *emptyIterator) HasNext() bool {
	return false
}

func (e *emptyIterator) Next() interface{} {
	return nil
}

func Empty() Iterable {
	return NewIterable(func() Iterator {
		return &emptyIterator{}
	})
}
