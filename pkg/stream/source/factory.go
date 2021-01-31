package source

import (
	"github.com/nicholasham/piper/pkg/core"
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
	return FromIterable(core.Range(start, end))
}

func FromIterable(iterable core.Iterable) *stream.SourceGraph {
	return Single(iterable).
		MapConcat(toIterable).
		Named("IterableSource")
}

func toIterable(value interface{}) (core.Iterable, error) {
	return value.(core.Iterable), nil
}
