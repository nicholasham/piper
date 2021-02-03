package source

import (
	"github.com/nicholasham/piper/pkg/core/iterable"
	"github.com/nicholasham/piper/pkg/stream"
)

func Concat(graphs ...*stream.SourceGraph) *stream.SourceGraph {
	return stream.ConcatSources(graphs...)
}

func Interleave(segmentSize int, graphs ...*stream.SourceGraph) *stream.SourceGraph {
	return stream.InterleaveSources(segmentSize, graphs...)
}

func Merge(graphs ...*stream.SourceGraph) *stream.SourceGraph {
	return stream.MergeSources(graphs...)
}

func Single(value interface{}) *stream.SourceGraph {
	return stream.FromSource(stream.SingleSource(value))
}

func Range(start int, end int) *stream.SourceGraph {
	return FromIterable(iterable.Range(start, end))
}

func FromIterable(iterable iterable.Iterable) *stream.SourceGraph {
	return Single(iterable).
		MapConcat(toIterable).
		Named("IterableSource")
}

func toIterable(value interface{}) (iterable.Iterable, error) {
	return value.(iterable.Iterable), nil
}
