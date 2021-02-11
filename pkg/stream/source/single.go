package source

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/stream"
)

// verify singleSourceStage implements stream.SourceStage interface
var _ stream.SourceStage = (*singleSourceStage)(nil)

type singleSourceStage struct {
	attributes *stream.StageAttributes
	value      interface{}
}

func (s *singleSourceStage) Open(_ context.Context, _ stream.MaterializeFunc) (stream.Reader, *core.Future) {
	outputPromise := core.NewPromise()
	outputStream := stream.NewStream()
	go func() {
		writer := outputStream.Writer()
		defer writer.Close()
		writer.SendValue(s.value)
		outputPromise.TrySuccess(stream.NotUsed)
	}()
	return outputStream.Reader(), outputPromise.Future()
}

func (s *singleSourceStage) With(options ...stream.StageOption) stream.Stage {
	return &singleSourceStage{
		attributes: s.attributes.Apply(options...),
		value:      s.value,
	}
}

func singleStage(value interface{}) stream.SourceStage {
	return &singleSourceStage{
		attributes: stream.DefaultStageAttributes,
		value:      value,
	}
}
