package experiment

import (
	"context"
	. "github.com/nicholasham/piper/pkg/core"
)

type Stage interface {
	With(options ...StageOption) Stage
}

type SourceStage interface {
	Stage
	UpstreamStage
}

type FlowStage interface {
	Stage
	UpstreamStage
	WireTo(stage UpstreamStage) FlowStage
}

type SinkStage interface {
	Stage
	WireTo(stage UpstreamStage) SinkStage
	Run(ctx context.Context, mat MaterializeFunc) *Promise
}

type UpstreamStage interface {
	Open(ctx context.Context, mat MaterializeFunc) (StreamReader, *Promise)
}

type MaterializeFunc func (left *Promise, right *Promise) *Promise
type MapMaterializedValueFunc func(value interface{}) interface{}

func KeepLeft(left *Promise, right *Promise) *Promise {
	return left
}

func KeepRight(left *Promise, right *Promise) *Promise {
	return right
}

func KeepBoth(left *Promise, right *Promise) *Promise {
	return left
}

type Future interface {
	Await() Result
}
