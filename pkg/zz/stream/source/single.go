package source

import (
	"context"
	"github.com/nicholasham/piper/pkg/streamold"
)

// verify iteratorSource implements stream.SourceStage interface
var _ streamold.SourceStage = (*singleSourceStage)(nil)

type singleSourceStage struct {
	name    string
	element streamold.Element
	outlet  *streamold.Outlet
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

func (s *singleSourceStage) Outlet() *streamold.Outlet {
	return s.outlet
}


func singleStage(value interface{}, options ...streamold.StageOption) streamold.SourceStage {

	stageOptions := streamold.DefaultStageOptions.
		Apply(streamold.Name("SingleSource")).
		Apply(options...)

	return &singleSourceStage{
		name:    stageOptions.Name,
		element: streamold.Value(value),
		outlet:  streamold.NewOutlet(stageOptions),
	}
}

