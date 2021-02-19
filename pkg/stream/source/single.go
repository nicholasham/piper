package source

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/stream"
	"sync"
)

// verify singleSourceStage implements stream.SourceStage interface
var _ stream.SourceStage = (*singleSourceStage)(nil)

type singleSourceStage struct {
	attributes *stream.StageAttributes
	value      interface{}
}

func (s *singleSourceStage) Named(name string) stream.Stage {
	return s.With(stream.Name(name))
}

func (s *singleSourceStage) Open(_ context.Context, wg *sync.WaitGroup, _ stream.MaterializeFunc) (*stream.Receiver, *core.Future) {
	outputPromise := core.NewPromise()
	outputStream := stream.NewStream(s.attributes.Name)
	wg.Add(1)
	go func() {
		writer := outputStream.Sender()
		defer func() {
			writer.Close()
			wg.Done()
		}()
		writer.TrySend(stream.Value(s.value))
		outputPromise.TrySuccess(stream.NotUsed)
	}()
	return outputStream.Receiver(), outputPromise.Future()
}

func (s *singleSourceStage) With(options ...stream.StageOption) stream.Stage {
	return &singleSourceStage{
		attributes: s.attributes.With(options...),
		value:      s.value,
	}
}

func singleStage(value interface{}) stream.SourceStage {
	return &singleSourceStage{
		attributes: stream.DefaultStageAttributes.With(stream.Name("SingleSource")),
		value:      value,
	}
}
