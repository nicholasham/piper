package source

import (
	"github.com/nicholasham/piper/pkg/piper"
	"github.com/nicholasham/piper/pkg/types"
	"github.com/nicholasham/piper/pkg/types/iterator"
)

// Range Emit each integer in a stepped range.
func Range(start int, end int, step int, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return Iterator(iterator.Range(start, end, step), append(attributes, piper.Name("RangeSource"))...)
}

func Iterator(iterator iterator.Iterator, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return piper.SourceFrom(iteratorSource(iterator, attributes))
}

func List(values []interface{}, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return Iterator(iterator.Slice(values...), append(attributes, piper.Name("ListSource"))...)
}

func Failed(cause error, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return piper.SourceFrom(failedSource(cause, attributes))
}

func Empty(attributes ...piper.StageAttribute) *piper.SourceGraph {
	return piper.SourceFrom(emptySource(attributes))
}

func Unfold(state interface{}, f iterator.UnfoldFunc, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return Iterator(iterator.Unfold(state, f), attributes...)
}

func Repeat(value interface{}, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return Unfold(value, func(state interface{}) types.Option {
		return types.Some(value)
	}, attributes...)
}

func Concat(graphs ...*piper.SourceGraph) piper.SourceGraphFactory {
	return func(attributes ...piper.StageAttribute) *piper.SourceGraph {
		return piper.CombineSources(graphs)(piper.ConcatStrategy(), append(attributes, piper.Name("ConcatSource"))...)
	}
}

func Merge(graphs ...*piper.SourceGraph) piper.SourceGraphFactory {
	return func(attributes ...piper.StageAttribute) *piper.SourceGraph {
		return piper.CombineSources(graphs)(piper.MergeStrategy(), append(attributes, piper.Name("MergeSource"))...)
	}
}

func Interleave(segmentSize int, graphs ...*piper.SourceGraph) piper.SourceGraphFactory {
	return func(attributes ...piper.StageAttribute) *piper.SourceGraph {
		return piper.CombineSources(graphs)(piper.InterleaveStrategy(segmentSize), append(attributes, piper.Name("InterleaveSource"))...)
	}
}
