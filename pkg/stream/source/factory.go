package source

import "github.com/nicholasham/piper/pkg/stream"

func Concat(graphs []*stream.SourceGraph) *stream.SourceGraph {
	return stream.CombineSources(graphs, stream.ConcatStrategy()).With(stream.Name("ConcatSource"))
}

func Interleave(segmentSize int, graphs []*stream.SourceGraph) *stream.SourceGraph {
	return stream.CombineSources(graphs, stream.InterleaveStrategy(segmentSize)).With(stream.Name("InterleaveSource"))
}

func Merge(graphs []*stream.SourceGraph) *stream.SourceGraph {
	return stream.CombineSources(graphs, stream.ConcatStrategy()).With(stream.Name("MergeSource"))
}

