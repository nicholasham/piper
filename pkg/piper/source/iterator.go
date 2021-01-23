package source

import (
	"context"
	"github.com/nicholasham/piper/pkg/piper"
	"github.com/nicholasham/piper/pkg/piper/attribute"
	"github.com/nicholasham/piper/pkg/types/iterator"
)

// verify iteratorSource implements piper.SourceStage interface
var _ piper.SourceStage = (*iteratorSourceStage)(nil)

type iteratorSourceStage struct {
	attributes *attribute.StageAttributes
	outlet     *piper.Outlet
	iterator   iterator.Iterator
}

func (receiver *iteratorSourceStage) Name() string {
	return receiver.attributes.Name
}

func (receiver *iteratorSourceStage) Run(ctx context.Context) {
	go func(outlet *piper.Outlet, iterator iterator.Iterator) {
		defer outlet.Close()
		for iterator.HasNext() {
			select {
			case <-ctx.Done():
				outlet.Send(piper.Error(ctx.Err()))
				return
			case <-outlet.Done():
				return
			default:
			}
			value, err := iterator.Next()
			if err != nil {
				outlet.Send(piper.Error(err))
			} else {
				outlet.Send(piper.Value(value))
			}

		}
	}(receiver.outlet, receiver.iterator)
}

func (receiver *iteratorSourceStage) Outlet() *piper.Outlet {
	return receiver.outlet
}

func iteratorSource(iterator iterator.Iterator, attributes []attribute.StageAttribute) piper.SourceStage {
	stageAttributes := attribute.Default("IteratorSource", attributes...)
	return &iteratorSourceStage{
		attributes: stageAttributes,
		outlet:     piper.NewOutlet(stageAttributes),
		iterator:   iterator,
	}
}
