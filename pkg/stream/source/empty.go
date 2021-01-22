package source

import (
	"context"
	"github.com/nicholasham/piper/pkg/stream"
	"github.com/nicholasham/piper/pkg/stream/attribute"
)

// verify emptySourceStage implements SourceStage interface
var _ stream.SourceStage = (*emptySourceStage)(nil)

type emptySourceStage struct {
	attributes *attribute.StageAttributes
	outlet     *stream.Outlet
}

func (receiver *emptySourceStage) Name() string {
	return receiver.attributes.Name
}

func (receiver *emptySourceStage) Run(ctx context.Context) {
	go func(outlet *stream.Outlet) {
		outlet.Send(stream.Errorf("first of empty stream."))
		outlet.Close()
	}(receiver.outlet)
}

func (receiver *emptySourceStage) Outlet() *stream.Outlet {
	return receiver.outlet
}

func emptySource(attributes []attribute.StageAttribute) stream.SourceStage {
	stageAttributes := attribute.Default("FailedSource", attributes...)
	return &emptySourceStage{
		attributes: stageAttributes,
		outlet:     stream.NewOutlet(stageAttributes)}
}
