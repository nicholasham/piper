package source

import "github.com/nicholasham/piper/pkg/stream"

func Concat(graphs ...*stream.SourceGraph) *stream.SourceGraph {
	return stream.ConcatSources(graphs...)
}

func Interleave(segmentSize int, graphs ...*stream.SourceGraph) *stream.SourceGraph {
	return stream.InterleaveSources(segmentSize, graphs...)
}

func Merge(graphs ...*stream.SourceGraph) *stream.SourceGraph {
	return stream.MergeSources(graphs...)
}
