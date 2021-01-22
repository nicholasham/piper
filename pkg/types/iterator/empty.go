package iterator

// verify emptyIterator implements Iterator interface
var _ Iterator = (*emptyIterator)(nil)

type emptyIterator struct {
}


func (e *emptyIterator) ToList() []T {
	return toList(e)
}

func (e *emptyIterator) HasNext() bool {
	return false
}

func (e *emptyIterator) Next() (T, error) {
	return nil, EmptyError
}

func Empty() Iterator {
	return &emptyIterator{}
}
