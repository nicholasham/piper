package experiment

import (
	"context"
	. "github.com/ahmetb/go-linq/v3"
	"github.com/nicholasham/piper/pkg/streamold"
	"github.com/nicholasham/piper/pkg/types"
)

type MapConcat func(interface{}) Query

type MapConcatFactory func() MapConcat

var _ streamold.SourceStage = (*selectManyStage)(nil)

type selectManyStage struct {
	name    string
	inlet   *streamold.Inlet
	outlet  *streamold.Outlet
	factory MapConcatFactory
	decider streamold.Decider
}

func (s *selectManyStage) Name() string {
	return s.name
}

func (s *selectManyStage) Run(ctx context.Context) {
	go func() {
		defer s.outlet.Close()
		f := s.factory()
		for element := range s.inlet.In() {

			select {
			case <-ctx.Done():
				s.outlet.Send(streamold.Error(ctx.Err()))
				return
			case <-s.outlet.Done():
				return
			default:
			}

			query := f(element.Value())
			next := query.Iterate()
			for item, ok := next(); ok; {
				result, ok := item.(types.Result)
				if !ok {
					s.inlet.Complete()
					break
				}

				result.IfSuccess(func(value interface{}) {
					s.outlet.Send(streamold.Value(value))
				})
				result.IfFailure(func(err error) {
					switch s.decider(err) {
					case streamold.Stop:
						s.inlet.Complete()
					case streamold.Resume:
					case streamold.Reset:
					}
				})
			}
		}
	}()
}

func (s *selectManyStage) Outlet() *streamold.Outlet {
	return s.outlet
}
