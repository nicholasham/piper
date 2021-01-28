package source

import (
	"context"

	"github.com/nicholasham/piper/pkg/zz/stream"
)

// verify failedSourceStage implements SourceStage interface
var _ stream.SourceStage = (*failedSourceStage)(nil)

type failedSourceStage struct {
	attributes *stream.StageOptions
	cause      error
	outlet     *stream.Outlet
}

func (receiver *failedSourceStage) Name() string {
	return receiver.attributes.Name
}

func (receiver *failedSourceStage) Run(ctx context.Context) {
	go func(outlet *stream.Outlet, cause error) {
		outlet.Send(stream.Error(cause))
		outlet.Close()
		return
	}(receiver.outlet, receiver.cause)
}

func (receiver *failedSourceStage) Outlet() *stream.Outlet {
	return receiver.outlet
}

func failedSource(cause error, options ...stream.StageOption) stream.SourceStage {
	stageOptions := stream.DefaultStageOptions.
						Apply(stream.Name("FailedSource")).
						Apply(options...)
	return &failedSourceStage{
		attributes: stageOptions,
		cause:      cause,
		outlet:     stream.NewOutlet(stageOptions),
	}
}
