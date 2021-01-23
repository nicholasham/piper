package source

import (
	"context"

	"github.com/nicholasham/piper/pkg/stream"
)

// verify failedSourceStage implements SourceStage interface
var _ stream.SourceStage = (*failedSourceStage)(nil)

type failedSourceStage struct {
	attributes *stream.StageAttributes
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

func failedSource(cause error, attributes []stream.StageAttribute) stream.SourceStage {
	stageAttributes := stream.NewAttributes("FailedSource", attributes...)
	return &failedSourceStage{
		attributes: stageAttributes,
		cause:      cause,
		outlet:     stream.NewOutlet(stageAttributes),
	}
}
