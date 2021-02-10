package stream

import (
	. "github.com/nicholasham/piper/pkg/core"
)

type MaterializeFunc func(left *Future, right *Future) *Future
type MapMaterializedValueFunc func(value Any) Result

func KeepLeft(left *Future, right *Future) *Future {
	return left
}

func KeepRight(left *Future, right *Future) *Future {
	return right
}

func KeepBoth(left *Future, right *Future) *Future {
	return left
}
