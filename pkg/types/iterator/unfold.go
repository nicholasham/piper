package iterator

import (
	"github.com/nicholasham/piper/pkg/types/optional"
)

type UnfoldFunc func(state interface{}) optional.Option

// verify unfoldIterator implements Iterator interface
var _ Iterator = (*unfoldIterator)(nil)

type unfoldIterator struct {
	result optional.Option
	f      UnfoldFunc
}

func (u *unfoldIterator) HasNext() bool {
	return u.result.IsSome()
}

func (u *unfoldIterator) Next() (T, error) {

	value, err := u.result.Get()

	if err != nil {
		return nil, err
	}

	u.result = u.f(value)

	return value, nil

}

func (u *unfoldIterator) ToList() []T {
	return toList(u)
}

func Unfold(state interface{}, f UnfoldFunc) Iterator {
	return &unfoldIterator{
		result: optional.Some(state),
		f:      f,
	}
}
