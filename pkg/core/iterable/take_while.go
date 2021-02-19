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
	if !t.headDefined {
		if t.tail.HasNext() {
			t.head = t.tail.Next()
			if t.p(t.head) {
				t.headDefined = true
			}else {
				t.tail = Empty().Iterator()
			}

			if t.headDefined {
				return true
			}
		}
		return false
	}

	return true
}

func (t *takeWhileIterator) Next() interface{} {
	if t.HasNext() {
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