package iterable


// verify filterIterator implements Iterator interface
var _ Iterator = (*takeIterator)(nil)

type takeIterator struct {
	number        int
	taken int
	iterator Iterator
}

func (t *takeIterator) HasNext() bool {
	if t.taken == t.number {
		t.iterator = Empty().Iterator()
	}
	return t.iterator.HasNext()
}

func (t *takeIterator) Next() interface{} {
	t.taken ++
	return t.iterator.Next()
}

func take(number int, iterator Iterator) Iterable {
	return NewIterable(func() Iterator {
		return &takeIterator{
			iterator: iterator,
			number:        number,
		}
	})
}
