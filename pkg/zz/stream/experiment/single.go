package experiment

import "github.com/nicholasham/piper/pkg/zz/stream"

var _ OutHandler = (*SingleSource)(nil)

type SingleSource struct {
	value interface{}
}

func (s SingleSource) OnPull(actions StageActions) {
	actions.Push(stream.Value(s.value))
	actions.CompleteStage()
}

func (s SingleSource) OnDownstreamFinish(actions StageActions) {
}
