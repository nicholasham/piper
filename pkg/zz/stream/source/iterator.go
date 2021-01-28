package source

import (
	"context"

	"github.com/nicholasham/piper/pkg/streamold"
	"github.com/nicholasham/piper/pkg/types/iterator"
)

// verify iteratorSource implements stream.SourceStage interface
var _ streamold.SourceStage = (*iteratorSourceStage)(nil)

type iteratorSourceStage struct {
	outlet   *streamold.Outlet
	iterator iterator.Iterator
	name     string
}

func (receiver *iteratorSourceStage) Name() string {
	return receiver.name
}

func (receiver *iteratorSourceStage) Run(ctx context.Context) {
	go func(outlet *streamold.Outlet, iterator iterator.Iterator) {
		defer outlet.Close()
		for iterator.HasNext() {
			select {
			case <-ctx.Done():
				outlet.Send(streamold.Error(ctx.Err()))
				return
			case <-outlet.Done():
				return
			default:
			}
			value, err := iterator.Next()
			if err != nil {
				outlet.Send(streamold.Error(err))
			} else {
				outlet.Send(streamold.Value(value))
			}

		}
	}(receiver.outlet, receiver.iterator)
}

func (receiver *iteratorSourceStage) Outlet() *streamold.Outlet {
	return receiver.outlet
}

func iteratorSource(name string, iterator iterator.Iterator, options ...streamold.StageOption) streamold.SourceStage {
	stageOptions := streamold.DefaultStageOptions.
		Apply(streamold.Name("IteratorSource")).
		Apply(options...)

	return &iteratorSourceStage{
		name:     stageOptions.Name,
		outlet:   streamold.NewOutlet(stageOptions),
		iterator: iterator,
	}
}
