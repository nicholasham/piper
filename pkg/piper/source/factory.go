package source

import (
	"github.com/nicholasham/piper/pkg/piper"
	"github.com/nicholasham/piper/pkg/piper/attribute"
	"github.com/nicholasham/piper/pkg/types"
	"github.com/nicholasham/piper/pkg/types/iterator"
)

// Range Emit each integer in a stepped range.
func Range(start int, end int, step int, attributes ...attribute.StageAttribute) *piper.SourceGraph {
	return Iterator(iterator.Range(start, end, step), append(attributes, attribute.Name("RangeSource"))...)
}

func Iterator(iterator iterator.Iterator, attributes ...attribute.StageAttribute) *piper.SourceGraph {
	return piper.SourceFrom(iteratorSource(iterator, attributes))
}

func List(values []interface{}, attributes ...attribute.StageAttribute) *piper.SourceGraph {
	return Iterator(iterator.Slice(values...), append(attributes, attribute.Name("ListSource"))...)
}

func Failed(cause error, attributes ...attribute.StageAttribute) *piper.SourceGraph {
	return piper.SourceFrom(failedSource(cause, attributes))
}

func Empty(attributes ...attribute.StageAttribute) *piper.SourceGraph {
	return piper.SourceFrom(emptySource(attributes))
}

func Unfold(state interface{}, f iterator.UnfoldFunc, attributes ...attribute.StageAttribute) *piper.SourceGraph {
	return Iterator(iterator.Unfold(state, f), attributes...)
}

func Repeat(value interface{}, attributes ...attribute.StageAttribute) *piper.SourceGraph {
	return Unfold(value, func(state interface{}) types.Option {
		return types.Some(value)
	}, attributes...)
}

func Concat(graphs ...*piper.SourceGraph) piper.SourceGraphFactory {
	return func(attributes ...attribute.StageAttribute) *piper.SourceGraph {
		return piper.CombineSources(graphs)(piper.ConcatStrategy(), append(attributes, attribute.Name("ConcatSource"))...)
	}
}

func Merge(graphs ...*piper.SourceGraph) piper.SourceGraphFactory {
	return func(attributes ...attribute.StageAttribute) *piper.SourceGraph {
		return piper.CombineSources(graphs)(piper.MergeStrategy(), append(attributes, attribute.Name("MergeSource"))...)
	}
}

func Interleave(segmentSize int, graphs ...*piper.SourceGraph) piper.SourceGraphFactory {
	return func(attributes ...attribute.StageAttribute) *piper.SourceGraph {
		return piper.CombineSources(graphs)(piper.InterleaveStrategy(segmentSize), append(attributes, attribute.Name("InterleaveSource"))...)
	}
}
