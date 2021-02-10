package experiment

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
)

// verify mapConcat implements SourceStage interface
var _ SourceStage = (*singleStage)(nil)

type singleStage struct {
	attributes *StageAttributes
	value      interface{}
}

func (s *singleStage) Open(ctx context.Context, mat MaterializeFunc) (StreamReader, *core.Future) {
	outputPromise := core.NewPromise()
	outputStream := NewStream()
	go func(){
		writer := outputStream.Writer()
		writer.SendValue(s.value)
		writer.Close()
		outputPromise.TrySuccess(NotUsed)
	}()
	return outputStream.Reader(), outputPromise.Future()

}

func (s *singleStage) Name() string {
	return s.attributes.Name
}

func (s *singleStage) With(options ...StageOption) Stage {
	attributes := s.attributes.Apply(options...)
	return &singleStage{
		attributes: attributes,
		value:      s.value,
	}
}

func SingleStage(value interface{}) SourceStage {
	return &singleStage{
		attributes: DefaultStageAttributes,
		value:      value,
	}
}
