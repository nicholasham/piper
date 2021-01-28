package experiment

import (
	"context"
	"github.com/nicholasham/piper/pkg/streamold"
)

var _ streamold.SourceStage = (*sourceStage)(nil)
var _ StageActions = (*sourceStage)(nil)

type sourceStage struct {
	name   string
	outlet *streamold.Outlet
	logger streamold.Logger
	logic  OutHandler
}

func (s *sourceStage) CompleteStage() {
	s.outlet.Close()
}

func (s *sourceStage) Push(element streamold.Element) {
	s.outlet.Send(element)
	s.outlet.Close()
}

func (s *sourceStage) FailStage(cause error) {

}

func (s *sourceStage) Name() string {
	return s.name
}

func (s *sourceStage) Run(ctx context.Context) {
	go func() {
		s.logic.OnPull(s)
		s.logic.OnDownstreamFinish(s)
	}()
}

func (s *sourceStage) Outlet() *streamold.Outlet {
	return s.outlet
}

