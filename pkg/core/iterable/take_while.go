package iterable

import "github.com/nicholasham/piper/pkg/core"

// verify takeWhileIterator implements Iterator interface
var _ Iterator = (*takeWhileIterator)(nil)

type takeWhileIterator struct {
	p           core.PredicateFunc
	tail        Iterator
	head        interface{}
	headDefined bool
}

func (t *takeWhileIterator) HasNext() bool {
	t.head = t.tail.Next()
	if t.head != nil && t.p(t.head) {
		t.headDefined = true
	}else {
		t.tail = Empty().Iterator()
	}
	return t.headDefined
}

func (t *takeWhileIterator) Next() interface{} {
	if t.headDefined && t.tail.HasNext() {
		t.headDefined = false
		return t.head
	}else {
		return Empty().Iterator().Next()
	}
}

func takeWhile(p core.PredicateFunc, iterator Iterator) Iterable {
	return NewIterable(func() Iterator {
		return &takeWhileIterator {
			tail:        iterator,
			p:           p,
			headDefined: false,
		}
	})
}