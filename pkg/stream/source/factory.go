package source

import (
	"github.com/nicholasham/piper/pkg/core/iterable"
	"github.com/nicholasham/piper/pkg/stream"
)

func Single(value interface{}) *stream.SourceGraph {
	return stream.FromSource(singleStage(value))
}

func Failed(err error) *stream.SourceGraph {
	return stream.FromSource(failedStage(err))
}

func Empty() *stream.SourceGraph {
	return FromIterable(iterable.Empty())
}

func List(values ...interface{}) *stream.SourceGraph {
	return FromIterable(iterable.Slice(values...))
}

func Range(start int, end int) *stream.SourceGraph {
	return FromIterable(iterable.Range(start, end))
}

func FromIterable(iterable iterable.Iterable) *stream.SourceGraph {
	return Single(iterable).
		MapConcat(toIterable)
}

func toIterable(value interface{}) (iterable.Iterable, error) {
	return value.(iterable.Iterable), nil
}

func Concat(graphs ...*stream.SourceGraph) *stream.SourceGraph {
	return stream.ConcatSources(graphs...)
}

func Interleave(segmentSize int, graphs ...*stream.SourceGraph) *stream.SourceGraph {
	return stream.InterleaveSources(segmentSize, graphs...)
}

func Merge(graphs ...*stream.SourceGraph) *stream.SourceGraph {
	return stream.MergeSources(graphs...)
}
