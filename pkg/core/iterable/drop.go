package iterable


// verify filterIterator implements Iterator interface
var _ Iterator = (*dropIterator)(nil)

type dropIterator struct {
	number        int
	dropped int
	iterator Iterator
}

func (t *dropIterator) HasNext() bool {
	return t.iterator.HasNext()
}

func (t *dropIterator) Next() interface{} {
	t.dropped ++
	return t.iterator.Next()
}

func drop(number int, iterator Iterator) Iterable {
	return NewIterable(func() Iterator {
		return &dropIterator{
			iterator: iterator,
			number:        number,
		}
	})
}
