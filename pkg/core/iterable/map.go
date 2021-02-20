package iterable

import "github.com/nicholasham/piper/pkg/core"

// verify filterIterator implements Iterator interface
var _ Iterator = (*mapIterator)(nil)

type mapIterator struct {
	f        core.MapFunc
	iterator Iterator
}

func (t *mapIterator) HasNext() bool {
	return t.iterator.HasNext()
}

func (t *mapIterator) Next() interface{} {
	return t.f(t.iterator.Next())
}

func mapping(f core.MapFunc, iterator Iterator) Iterable {
	return NewIterable(func() Iterator {
		return &mapIterator{
			iterator: iterator,
			f:        f,
		}
	})
}
