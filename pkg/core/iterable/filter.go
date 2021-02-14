package iterable

import "github.com/nicholasham/piper/pkg/core"

// verify filterIterator implements Iterator interface
var _ Iterator = (*filterIterator)(nil)

type filterIterator struct {
	p           core.PredicateFunc
	iterator        Iterator
	head        interface{}
	headDefined bool
}

func (t *filterIterator) HasNext() bool {
	for  {
		if !t.iterator.HasNext()  {
			return false
		}
		t.head =  t.iterator.Next()
		if t.p(t.head) {
			t.headDefined = true
			return true
		}
	}
}

func (t *filterIterator) Next() interface{} {
	if t.headDefined {
		t.headDefined = false
		return t.head
	}
	return Empty().Iterator().Next()
}

func filter(p core.PredicateFunc, iterator Iterator) Iterable {
	return NewIterable(func() Iterator {
		return &filterIterator {
			iterator:        iterator,
			p:           p,
			headDefined: false,
		}
	})
}