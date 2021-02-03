package iterator

import (
	"github.com/nicholasham/piper/pkg/core"
)

type UnfoldFunc func(state interface{}) core.Optional

// verify unfoldIterator implements Iterator interface
var _ Iterator = (*unfoldIterator)(nil)

type unfoldIterator struct {
	result core.Optional
	f      UnfoldFunc
}

func (u *unfoldIterator) HasNext() bool {
	return u.result.IsSome()
}

func (u *unfoldIterator) Next() (interface{}, error) {

	value, err := u.result.Get()

	if err != nil {
		return nil, err
	}

	u.result = u.f(value)

	return value, nil

}

func (u *unfoldIterator) ToList() []interface{} {
	return toList(u)
}

func Unfold(state interface{}, f UnfoldFunc) Iterator {
	return &unfoldIterator{
		result: core.Some(state),
		f:      f,
	}
}
