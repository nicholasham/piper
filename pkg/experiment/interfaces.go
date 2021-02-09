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
	Run(ctx context.Context, mat MaterializeFunc) *Future
}

type UpstreamStage interface {
	Open(ctx context.Context, mat MaterializeFunc) (StreamReader, *Future)
}

type MaterializeFunc func (left *Future, right *Future) *Future
type MapMaterializedValueFunc func(value Any) Any

func KeepLeft(left *Future, right *Future) *Future {
	return left
}

func KeepRight(left *Future, right *Future) *Future {
	return right
}

func KeepBoth(left *Future, right *Future) *Future {
	return left
}