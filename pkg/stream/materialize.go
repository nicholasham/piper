package stream

import (
	"github.com/nicholasham/piper/pkg/core"
)

type MaterializeFunc func(left *core.Future, right *core.Future) *core.Future
type MapMaterializedValueFunc func(value core.Any) core.Result

func KeepLeft(left *core.Future, right *core.Future) *core.Future {
	return left
}

func KeepRight(left *core.Future, right *core.Future) *core.Future {
	return right
}

func KeepBoth(left *core.Future, right *core.Future) *core.Future {
	return left
}
