package source

import (
	"context"
	"github.com/nicholasham/piper/pkg/stream"
	"github.com/nicholasham/piper/pkg/stream/attribute"
	"github.com/nicholasham/piper/pkg/types/iterator"
)

// verify iteratorSource implements stream.SourceStage interface
var _ stream.SourceStage = (*iteratorSourceStage)(nil)

type iteratorSourceStage struct {
	attributes *attribute.StageAttributes
	outlet     *stream.Outlet
	iterator   iterator.Iterator
}

func (receiver *iteratorSourceStage) Name() string {
	return receiver.attributes.Name
}

func (receiver *iteratorSourceStage) Run(ctx context.Context) {
	go func(outlet *stream.Outlet, iterator iterator.Iterator) {
		defer outlet.Close()
		for iterator.HasNext() {
			select {
			case <-ctx.Done():
				outlet.Send(stream.Error(ctx.Err()))
				return
			case <-outlet.Done():
				return
			default:
			}
			value, err := iterator.Next()
			if err != nil {
				outlet.Send(stream.Error(err))
			} else {
				outlet.Send(stream.Value(value))
			}

		}
	}(receiver.outlet, receiver.iterator)
}

func (receiver *iteratorSourceStage) Outlet() *stream.Outlet {
	return receiver.outlet
}

func iteratorSource(iterator iterator.Iterator, attributes []attribute.StageAttribute) stream.SourceStage {
	stageAttributes := attribute.Default("IteratorSource", attributes...)
	return &iteratorSourceStage{
		attributes: stageAttributes,
		outlet:     stream.NewOutlet(stageAttributes),
		iterator:   iterator,
	}
}
