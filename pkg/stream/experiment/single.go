package experiment

import "github.com/nicholasham/piper/pkg/stream"

var _ OutHandler = (*SingleSource)(nil)

type SingleSource struct {
	value interface{}
}

func (s SingleSource) OnPull(actions StageActions) {
	actions.Push(stream.Value(s.value))
}

func (s SingleSource) OnDownstreamFinish(actions StageActions) {
}


