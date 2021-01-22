package source

import (
	"github.com/nicholasham/piper/pkg/stream"
	"github.com/nicholasham/piper/pkg/stream/attribute"
	"github.com/nicholasham/piper/pkg/types"
	"github.com/nicholasham/piper/pkg/types/iterator"
)

// Range Emit each integer in a stepped range.
func Range(start int, end int, step int, attributes ...attribute.StageAttribute) *stream.SourceGraph {
	return Iterator(iterator.Range(start, end, step), append(attributes, attribute.Name("RangeSource"))...)
}

func Iterator(iterator iterator.Iterator, attributes ...attribute.StageAttribute) *stream.SourceGraph {
	return stream.SourceFrom(iteratorSource(iterator, attributes))
}

func List(values []types.T, attributes ...attribute.StageAttribute) *stream.SourceGraph {
	return Iterator(iterator.Slice(values), append(attributes, attribute.Name("ListSource"))...)
}

func Failed(cause error, attributes ...attribute.StageAttribute) *stream.SourceGraph {
	return stream.SourceFrom(failedSource(cause, attributes))
}

func Empty(attributes ...attribute.StageAttribute) *stream.SourceGraph {
	return stream.SourceFrom(emptySource(attributes))
}

func Unfold(state interface{}, f iterator.UnfoldFunc, attributes ...attribute.StageAttribute) *stream.SourceGraph {
	return Iterator(iterator.Unfold(state, f), attributes...)
}

func Repeat(value interface{}, attributes ...attribute.StageAttribute) *stream.SourceGraph {
	return Unfold(value, func(state interface{}) types.Option {
		return types.Some(value)
	}, attributes...)
}

func Concat(graphs ...*stream.SourceGraph) stream.SourceGraphFactory {
	return func(attributes ...attribute.StageAttribute) *stream.SourceGraph {
		return stream.CombineSources(graphs)(stream.ConcatStrategy(), append(attributes, attribute.Name("ConcatSource"))...)
	}
}

func Merge(graphs ...*stream.SourceGraph) stream.SourceGraphFactory {
	return func(attributes ...attribute.StageAttribute) *stream.SourceGraph {
		return stream.CombineSources(graphs)(stream.MergeStrategy(), append(attributes, attribute.Name("MergeSource"))...)
	}
}

func Interleave(segmentSize int, graphs ...*stream.SourceGraph) stream.SourceGraphFactory {
	return func(attributes ...attribute.StageAttribute) *stream.SourceGraph {
		return stream.CombineSources(graphs)(stream.InterleaveStrategy(segmentSize), append(attributes, attribute.Name("InterleaveSource"))...)
	}
}
