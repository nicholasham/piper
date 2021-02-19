package source

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/stream"
	"sync"
)

// verify failedSourceStage implements stream.SourceStage interface
var _ stream.SourceStage = (*failedSourceStage)(nil)


type failedSourceStage struct {
	attributes *stream.StageAttributes
	err error
}

func (f *failedSourceStage) Named(name string) stream.Stage {
	return f.With(stream.Name(name))
}

func (f *failedSourceStage) With(options ...stream.StageOption) stream.Stage {
	return &failedSourceStage{
		attributes: f.attributes.With(options...),
	}
}

func (f *failedSourceStage) Open(_ context.Context, wg *sync.WaitGroup, _ stream.MaterializeFunc) (*stream.Receiver, *core.Future) {
	outputPromise := core.NewPromise()
	outputStream := stream.NewStream(f.attributes.Name)
	wg.Add(1)
	go func() {
		writer := outputStream.Sender()
		defer func() {
			writer.Close()
			wg.Done()
		}()
		writer.TrySend(stream.Error(f.err))
		outputPromise.TryFailure(f.err)
	}()
	return outputStream.Receiver(), outputPromise.Future()
}


func failedStage(err error) stream.SourceStage {
	return &failedSourceStage{
		attributes: stream.DefaultStageAttributes,
		err:      err,
	}
}
