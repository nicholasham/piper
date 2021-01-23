package source

import (
	"github.com/nicholasham/piper/pkg/stream"
	"github.com/nicholasham/piper/pkg/types"
	"github.com/nicholasham/piper/pkg/types/iterator"
)

// Range Emit each integer in a stepped range.
func Range(start int, end int, step int, attributes ...stream.StageAttribute) *stream.SourceGraph {
	return Iterator("RangeSource", iterator.Range(start, end, step), attributes...)
}

func Iterator(name string, iterator iterator.Iterator, attributes ...stream.StageAttribute) *stream.SourceGraph {
	return stream.SourceFrom(iteratorSource(name, iterator, attributes))
}

func List(values []interface{}, attributes ...stream.StageAttribute) *stream.SourceGraph {
	return Iterator("ListSource", iterator.Slice(values...), attributes...)
}

func Failed(cause error, attributes ...stream.StageAttribute) *stream.SourceGraph {
	return stream.SourceFrom(failedSource(cause, attributes))
}

func Empty(attributes ...stream.StageAttribute) *stream.SourceGraph {
	return Iterator("EmptySource", iterator.Empty(), attributes...)
}

func Unfold(state interface{}, f iterator.UnfoldFunc, attributes ...stream.StageAttribute) *stream.SourceGraph {
	return Iterator("UnfoldSource", iterator.Unfold(state, f), attributes...)
}

func Repeat(value interface{}, attributes ...stream.StageAttribute) *stream.SourceGraph {
	f := func(state interface{}) types.Option {
		return types.Some(value)
	}
	return Iterator("RepeatSource", iterator.Unfold(value, f), attributes...)
}

func Concat(graphs []*stream.SourceGraph, attributes ...stream.StageAttribute) *stream.SourceGraph {
	return stream.CombineSources("ConcatSource", graphs, stream.ConcatStrategy(), attributes...)

}

func Merge(graphs []*stream.SourceGraph, attributes ...stream.StageAttribute) *stream.SourceGraph {
	return stream.CombineSources("MergeSource", graphs, stream.MergeStrategy(), attributes...)

}

func Interleave(segmentSize int, graphs []*stream.SourceGraph, attributes ...stream.StageAttribute) *stream.SourceGraph {
	return stream.CombineSources("InterleaveSource", graphs, stream.InterleaveStrategy(segmentSize), attributes...)

}
