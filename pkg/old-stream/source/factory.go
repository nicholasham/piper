package source

import (
	"github.com/nicholasham/piper/pkg/core/iterable"
	"github.com/nicholasham/piper/pkg/old-stream"
)

func Concat(graphs ...*old_stream.SourceGraph) *old_stream.SourceGraph {
	return old_stream.ConcatSources(graphs...)
}

func Interleave(segmentSize int, graphs ...*old_stream.SourceGraph) *old_stream.SourceGraph {
	return old_stream.InterleaveSources(segmentSize, graphs...)
}

func Merge(graphs ...*old_stream.SourceGraph) *old_stream.SourceGraph {
	return old_stream.MergeSources(graphs...)
}

func Single(value interface{}) *old_stream.SourceGraph {
	return old_stream.FromSource(old_stream.SingleSource(value))
}

func Range(start int, end int) *old_stream.SourceGraph {
	return FromIterable(iterable.Range(start, end))
}

func FromIterable(iterable iterable.Iterable) *old_stream.SourceGraph {
	return Single(iterable).
		MapConcat(toIterable).
		Named("IterableSource")
}

func toIterable(value interface{}) (iterable.Iterable, error) {
	return value.(iterable.Iterable), nil
}
