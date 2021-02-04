package experiment

import (
	"github.com/nicholasham/piper/pkg/core"
	"golang.org/x/net/context"
)

type UpstreamStage interface {
	Open(ctx context.Context, mat MaterializeFunc) (StreamReader, *core.Promise)
}


type Inlet interface {
	Complete()
}

type Outlet interface {
	Close()
	SendError(value interface{})
	SendValue(value interface{})
}

type Element struct {

}

type MaterializeFunc func (left *core.Promise, right *core.Promise) *core.Promise

func KeepLeft(left *core.Promise, right *core.Promise) *core.Promise {
	return left
}

func KeepRight(left *core.Promise, right *core.Promise) *core.Promise {
	return left
}

func KeepBoth(left *core.Promise, right *core.Promise) *core.Promise {
	return left
}
