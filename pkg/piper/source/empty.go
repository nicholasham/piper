package source

import (
	"context"
	"github.com/nicholasham/piper/pkg/piper"
	"github.com/nicholasham/piper/pkg/piper/attribute"
)

// verify emptySourceStage implements SourceStage interface
var _ piper.SourceStage = (*emptySourceStage)(nil)

type emptySourceStage struct {
	attributes *attribute.StageAttributes
	outlet     *piper.Outlet
}

func (receiver *emptySourceStage) Name() string {
	return receiver.attributes.Name
}

func (receiver *emptySourceStage) Run(ctx context.Context) {
	go func(outlet *piper.Outlet) {
		outlet.Send(piper.Errorf("first of empty piper."))
		outlet.Close()
	}(receiver.outlet)
}

func (receiver *emptySourceStage) Outlet() *piper.Outlet {
	return receiver.outlet
}

func emptySource(attributes []attribute.StageAttribute) piper.SourceStage {
	stageAttributes := attribute.Default("FailedSource", attributes...)
	return &emptySourceStage{
		attributes: stageAttributes,
		outlet:     piper.NewOutlet(stageAttributes)}
}
