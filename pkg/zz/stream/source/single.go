package source

import (
	"context"

	"github.com/nicholasham/piper/pkg/zz/stream"
)

// verify iteratorSource implements stream.SourceStage interface
var _ stream.SourceStage = (*singleSourceStage)(nil)

type singleSourceStage struct {
	name    string
	element stream.Element
	outlet  *stream.Outlet
}

func (s *singleSourceStage) Name() string {
	return s.name
}

func (s *singleSourceStage) Run(ctx context.Context) {
	go func() {
		s.outlet.Send(s.element)
		s.outlet.Close()
	}()
}

func (s *singleSourceStage) Outlet() *stream.Outlet {
	return s.outlet
}

func singleStage(value interface{}, options ...stream.StageOption) stream.SourceStage {

	stageOptions := stream.DefaultStageOptions.
		Apply(stream.Name("SingleSource")).
		Apply(options...)

	return &singleSourceStage{
		name:    stageOptions.Name,
		element: stream.Value(value),
		outlet:  stream.NewOutlet(stageOptions),
	}
}
