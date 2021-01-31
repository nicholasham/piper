package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
)

// verify mapConcat implements FlowStage interface
var _ FlowStage = (*mapConcat)(nil)

type MapConcatFunc func(value interface{}) (core.Iterable, error)

type mapConcat struct {
	attributes *StageAttributes
	inlet      *Inlet
	outlet     *Outlet
	f          MapConcatFunc
}

func (s *mapConcat) Name() string {
	return s.attributes.Name
}

func (s *mapConcat) Run(ctx context.Context) {
	go func() {
		defer s.outlet.Close()
		for element := range s.inlet.In() {

			select {
			case <-ctx.Done():
				s.outlet.SendError(ctx.Err())
				return
			case <-s.outlet.Done():
				return
			default:
			}

			element.
				WhenError(s.outlet.SendError).
				WhenValue(s.handleValue)
		}
	}()
}

func (s *mapConcat) handleError(err error) {
	if err != nil {
		switch s.attributes.Decider(err) {
		case Stop:
			s.inlet.Complete()
		case Resume:
			s.outlet.SendError(err)
		case Reset:
		}
	}
}

func (s *mapConcat) handleValue(value interface{}) {
	iterable, err := s.f(value)

	if err != nil {
		s.handleError(err)
		return
	}

	iterable.ForEach(func(i interface{}) {
		s.Outlet().Send(Value(i))
	})
}

func (s *mapConcat) With(options ...StageOption) Stage {
	attributes := s.attributes.Apply(options...)
	return &mapConcat{
		attributes: attributes,
		inlet:      s.inlet,
		outlet:     NewOutlet(attributes),
		f:          s.f,
	}
	return s
}

func (s *mapConcat) WireTo(stage OutputStage) FlowStage {
	s.inlet.WireTo(stage.Outlet())
	return s
}

func (s *mapConcat) Outlet() *Outlet {
	return s.outlet
}

func MapConcatStage(f MapConcatFunc) FlowStage {
	attributes := DefaultStageAttributes.Apply(Name("MapConcatStage"))
	return &mapConcat{
		attributes: attributes,
		inlet:      NewInlet(attributes),
		outlet:     NewOutlet(attributes),
		f:          f,
	}
}
