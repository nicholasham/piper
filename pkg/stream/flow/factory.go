package source

import "github.com/nicholasham/piper/pkg/stream"

func Concat(graphs ...*stream.FlowGraph) *stream.FlowGraph {
	return stream.ConcatFlows(graphs...)
}

func Interleave(segmentSize int, graphs ...*stream.FlowGraph) *stream.FlowGraph {
	return stream.InterleaveFlows(segmentSize, graphs...)
}

func Merge(graphs ...*stream.FlowGraph) *stream.FlowGraph {
	return stream.MergeFlows(graphs...)
}
