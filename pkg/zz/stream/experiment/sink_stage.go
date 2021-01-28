package experiment

import (
	"context"
	"github.com/nicholasham/piper/pkg/streamold"
)

var _ streamold.SinkStage = (*sinkStage)(nil)

type sinkStage struct {
	name string
	inlet *streamold.Inlet
	promise *streamold.Promise
	logic InHandler
}

func (s *sinkStage) Name() string {
	return s.name
}

func (s *sinkStage) Run(ctx context.Context) {
	go func() {
		for element := range s.inlet.In() {
			select {
			case <-ctx.Done():
				s.promise.Reject(ctx.Err())
				s.inlet.Complete()
			default:
			}
			if !s.inlet.CompletionSignaled() {
				element.WhenError(func(err error) {
					s.logic.OnUpstreamFailure(err, s.newStageActions())
				}).WhenValue(func(value interface{}) {
					s.logic.OnPush(element)
				})
			}
		}
		s.logic.OnUpstreamFinish(s.newStageActions())
	}()
}

func (s *sinkStage) Wire(stage streamold.SourceStage) {
	s.Wire(stage)
}

func (s *sinkStage) Inlet() *streamold.Inlet {
	return s.inlet
}

func (s *sinkStage) Result() streamold.Future {
	return s.promise
}

func (s *sinkStage) newStageActions() StageActions {
	return & stageActions{
		onCompleteStage: func() {
			s.inlet.Complete()
		},
		onPush: func(element streamold.Element) {
		},
		onFailStage: func(cause error) {

		},
	}
}

