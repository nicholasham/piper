package source

import (
	"github.com/nicholasham/piper/pkg/piper"
	"github.com/nicholasham/piper/pkg/types"
	"github.com/nicholasham/piper/pkg/types/iterator"
)

// Range Emit each integer in a stepped range.
func Range(start int, end int, step int, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return Iterator("RangeSource", iterator.Range(start, end, step), attributes...)
}

func Iterator(name string, iterator iterator.Iterator, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return piper.SourceFrom(iteratorSource(name, iterator, attributes))
}

func List(values []interface{}, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return Iterator("ListSource", iterator.Slice(values...), attributes...)
}

func Failed(cause error, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return piper.SourceFrom(failedSource(cause, attributes))
}

func Empty(attributes ...piper.StageAttribute) *piper.SourceGraph {
	return Iterator("EmptySource", iterator.Empty(), attributes...)
}

func Unfold(state interface{}, f iterator.UnfoldFunc, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return Iterator("UnfoldSource", iterator.Unfold(state, f), attributes...)
}

func Repeat(value interface{}, attributes ...piper.StageAttribute) *piper.SourceGraph {
	f := func(state interface{}) types.Option {
		return types.Some(value)
	}
	return Iterator("RepeatSource", iterator.Unfold(value, f), attributes...)
}

func Concat(graphs []*piper.SourceGraph, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return piper.CombineSources("ConcatSource", graphs, piper.ConcatStrategy(), attributes...)

}

func Merge(graphs []*piper.SourceGraph, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return piper.CombineSources("MergeSource", graphs, piper.MergeStrategy(), attributes...)

}

func Interleave(segmentSize int, graphs []*piper.SourceGraph, attributes ...piper.StageAttribute) *piper.SourceGraph {
	return piper.CombineSources("InterleaveSource", graphs, piper.InterleaveStrategy(segmentSize), attributes...)

}
