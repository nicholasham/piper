package stream

import (
	"context"
)

// verify mapConcat implements SourceStage interface
var _ SourceStage = (*singleStage)(nil)

type singleStage struct {
	attributes *StageAttributes
	outlet     *Outlet
	value      interface{}
}

func (s *singleStage) Name() string {
	return s.attributes.Name
}

func (s *singleStage) Run(_ context.Context) {
	go func() {
		s.outlet.SendValue(s.value)
		s.outlet.Close()
	}()
}

func (s *singleStage) With(options ...StageOption) Stage {
	attributes := s.attributes.Apply(options...)
	return &singleStage{
		attributes: attributes,
		outlet:     NewOutlet(attributes),
		value:      s.value,
	}
}

func (s *singleStage) Outlet() *Outlet {
	return s.outlet
}

func SingleSource(value interface{}) SourceStage {
	attributes := DefaultStageAttributes.Apply(Name("SingleSource"))
	return &singleStage{
		attributes: attributes,
		outlet:     NewOutlet(attributes),
		value:      value,
	}
}
