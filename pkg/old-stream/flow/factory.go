package source

import "github.com/nicholasham/piper/pkg/old-stream"

func Concat(graphs ...*old_stream.FlowGraph) *old_stream.FlowGraph {
	return old_stream.ConcatFlows(graphs...)
}

func Interleave(segmentSize int, graphs ...*old_stream.FlowGraph) *old_stream.FlowGraph {
	return old_stream.InterleaveFlows(segmentSize, graphs...)
}

func Merge(graphs ...*old_stream.FlowGraph) *old_stream.FlowGraph {
	return old_stream.MergeFlows(graphs...)
}
