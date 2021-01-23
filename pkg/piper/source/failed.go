package source

import (
	"context"
	"github.com/nicholasham/piper/pkg/piper"
)

// verify failedSourceStage implements SourceStage interface
var _ piper.SourceStage = (*failedSourceStage)(nil)

type failedSourceStage struct {
	attributes *piper.StageAttributes
	cause      error
	outlet     *piper.Outlet
}

func (receiver *failedSourceStage) Name() string {
	return receiver.attributes.Name
}

func (receiver *failedSourceStage) Run(ctx context.Context) {
	go func(outlet *piper.Outlet, cause error) {
		outlet.Send(piper.Error(cause))
		outlet.Close()
		return
	}(receiver.outlet, receiver.cause)
}

func (receiver *failedSourceStage) Outlet() *piper.Outlet {
	return receiver.outlet
}

func failedSource(cause error, attributes []piper.StageAttribute) piper.SourceStage {
	stageAttributes := piper.NewAttributes("FailedSource", attributes...)
	return &failedSourceStage{
		attributes: stageAttributes,
		cause:      cause,
		outlet:     piper.NewOutlet(stageAttributes),
	}
}
