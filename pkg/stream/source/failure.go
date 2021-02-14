package source

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/stream"
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

func (f *failedSourceStage) Open(_ context.Context, _ stream.MaterializeFunc) (stream.Reader, *core.Future) {
	outputPromise := core.NewPromise()
	outputStream := stream.NewStream()
	go func() {
		writer := outputStream.Writer()
		writer.Send(stream.Error(f.err))
		outputPromise.TryFailure(f.err)
	}()
	return outputStream.Reader(), outputPromise.Future()
}


func failedStage(err error) stream.SourceStage {
	return &failedSourceStage{
		attributes: stream.DefaultStageAttributes,
		err:      err,
	}
}
