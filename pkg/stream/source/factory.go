package source

import (
	"github.com/ahmetb/go-linq/v3"
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

func FromQuery(query linq.Query) *stream.SourceGraph {
	return stream.FromSource(stream.SingleSource(query)).
		SelectMany(func(value interface{}) linq.Query {
		return value.(linq.Query)
	}).Named("QuerySource")
}

