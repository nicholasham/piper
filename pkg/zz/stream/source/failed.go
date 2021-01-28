package source

import (
	"context"

	"github.com/nicholasham/piper/pkg/streamold"
)

// verify failedSourceStage implements SourceStage interface
var _ streamold.SourceStage = (*failedSourceStage)(nil)

type failedSourceStage struct {
	attributes *streamold.StageOptions
	cause      error
	outlet     *streamold.Outlet
}

func (receiver *failedSourceStage) Name() string {
	return receiver.attributes.Name
}

func (receiver *failedSourceStage) Run(ctx context.Context) {
	go func(outlet *streamold.Outlet, cause error) {
		outlet.Send(streamold.Error(cause))
		outlet.Close()
		return
	}(receiver.outlet, receiver.cause)
}

func (receiver *failedSourceStage) Outlet() *streamold.Outlet {
	return receiver.outlet
}

func failedSource(cause error, options ...streamold.StageOption) streamold.SourceStage {
	stageOptions := streamold.DefaultStageOptions.
						Apply(streamold.Name("FailedSource")).
						Apply(options...)
	return &failedSourceStage{
		attributes: stageOptions,
		cause:      cause,
		outlet:     streamold.NewOutlet(stageOptions),
	}
}
