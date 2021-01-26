package experiment

import (
	"context"
	"github.com/nicholasham/piper/pkg/stream"
)

var _ stream.SourceStage = (*sourceStage)(nil)

type sourceStage struct {
	name string
	outlet *stream.Outlet
	logger stream.Logger
	logic OutHandler
}

func (s *sourceStage) Name() string {
	return s.name
}

func (s *sourceStage) Run(ctx context.Context) {
	go func() {
		s.logic.OnPull(s.newStageActions())
		s.logic.OnDownstreamFinish(s.newStageActions())
	}()
}

func (s *sourceStage) Outlet() *stream.Outlet {
	return s.outlet
}

func (s *sourceStage) newStageActions() StageActions {
	return & stageActions{
		onCompleteStage: func() {
			s.outlet.Close()
		},
		onPush: func(element stream.Element) {
			s.Outlet().Send(element)
		},
		onFailStage: func(cause error) {

		},
	}
}


