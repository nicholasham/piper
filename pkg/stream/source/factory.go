package source

import (
	"github.com/nicholasham/piper/pkg/stream"
	"github.com/nicholasham/piper/pkg/types"
	"github.com/nicholasham/piper/pkg/types/iterator"
)

// Range Emit each integer in a stepped range.
func Range(start int, end int, step int, options ...stream.StageOption) *stream.SourceGraph {
	return FromIterator("RangeSource", iterator.Range(start, end, step), options...)
}

func FromIterator(name string, iterator iterator.Iterator, options ...stream.StageOption) *stream.SourceGraph {
	return stream.SourceFrom(iteratorSource(name, iterator, options...))
}

func List(values []interface{}, options ...stream.StageOption) *stream.SourceGraph {
	return FromIterator("ListSource", iterator.Slice(values...), options...)
}

func Single(value interface{}, options ...stream.StageOption) *stream.SourceGraph {
	return stream.SourceFrom (singleStage(value, options...))
}

func Failed(cause error, options ...stream.StageOption) *stream.SourceGraph {
	return stream.SourceFrom(failedSource(cause, options...))
}

func Empty(options ...stream.StageOption) *stream.SourceGraph {
	return FromIterator("EmptySource", iterator.Empty(), options...)
}

func Unfold(state interface{}, f iterator.UnfoldFunc, options ...stream.StageOption) *stream.SourceGraph {
	return FromIterator("UnfoldSource", iterator.Unfold(state, f), options...)
}

func Repeat(value interface{}, options ...stream.StageOption) *stream.SourceGraph {
	f := func(state interface{}) types.Option {
		return types.Some(value)
	}
	return FromIterator("RepeatSource", iterator.Unfold(value, f), options...)
}

func Concat(graphs []*stream.SourceGraph, options ...stream.StageOption) *stream.SourceGraph {
	return stream.CombineSources("ConcatSource", graphs, stream.ConcatStrategy(), options...)

}

func Merge(graphs []*stream.SourceGraph, options ...stream.StageOption) *stream.SourceGraph {
	return stream.CombineSources("MergeSource", graphs, stream.MergeStrategy(), options...)

}

func Interleave(segmentSize int, graphs []*stream.SourceGraph, options ...stream.StageOption) *stream.SourceGraph {
	return stream.CombineSources("InterleaveSource", graphs, stream.InterleaveStrategy(segmentSize), options...)
}

func Combine(name string, graphs []*stream.SourceGraph, strategy stream.FanInStrategy, options ...stream.StageOption) *stream.SourceGraph {
	return stream.CombineSources(name, graphs, strategy, options...)
}
