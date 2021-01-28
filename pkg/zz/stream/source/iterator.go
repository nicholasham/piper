package source

import (
	"context"

	"github.com/nicholasham/piper/pkg/zz/stream"
	"github.com/nicholasham/piper/pkg/types/iterator"
)

// verify iteratorSource implements stream.SourceStage interface
var _ stream.SourceStage = (*iteratorSourceStage)(nil)

type iteratorSourceStage struct {
	outlet   *stream.Outlet
	iterator iterator.Iterator
	name     string
}

func (receiver *iteratorSourceStage) Name() string {
	return receiver.name
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

func iteratorSource(name string, iterator iterator.Iterator, options ...stream.StageOption) stream.SourceStage {
	stageOptions := stream.DefaultStageOptions.
		Apply(stream.Name("IteratorSource")).
		Apply(options...)

	return &iteratorSourceStage{
		name:     stageOptions.Name,
		outlet:   stream.NewOutlet(stageOptions),
		iterator: iterator,
	}
}
