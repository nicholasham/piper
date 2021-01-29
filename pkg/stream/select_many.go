package stream

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
)

// verify selectMany implements FlowStage interface
var _ FlowStage = (*selectMany)(nil)

type SelectManyFunc func(value interface{}) linq.Query
type selectMany struct {
	attributes *StageAttributes
	inlet *Inlet
	outlet *Outlet
	f SelectManyFunc
}

func (s *selectMany) Name() string {
	return s.attributes.Name
}

func (s *selectMany) Run(ctx context.Context) {
	f := s.f
	go func() {
		for element := range s.inlet.In() {
			element.WhenValue(func(value interface{}) {
				query := f(value)
				query.ForEach(func(i interface{}) {
					s.Outlet().Send(Value(i))
				})
			})
		}
	}()
}

func (s *selectMany) With(options ...StageOption) Stage {
	attributes := s.attributes.Apply(options...)
	return & selectMany{
		attributes: attributes,
		inlet: s.inlet,
		outlet: NewOutlet(attributes),
		f: s.f,
	}
	return s
}

func (s *selectMany) WireTo(stage OutputStage) FlowStage {
	s.inlet.WireTo(stage.Outlet())
	return s
}

func (s *selectMany) Outlet() *Outlet {
	return s.outlet
}

func SelectManyFlow(f SelectManyFunc) FlowStage {
	attributes := DefaultStageAttributes.Apply(Name("SelectManyFlow"))
	return & selectMany{
		attributes: attributes,
		inlet: NewInlet(attributes),
		outlet: NewOutlet(attributes),
		f: f,
	}
}
