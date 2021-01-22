package flow

import (
	"github.com/nicholasham/piper/pkg/stream"
	"github.com/nicholasham/piper/pkg/stream/attribute"
)

func Log(name string, attributes ...attribute.StageAttribute) *stream.FlowGraph {
	return stream.FlowFrom(logFlow(name, attributes...))
}

func Map(f MapFunc, attributes ...attribute.StageAttribute) *stream.FlowGraph {
	return stream.FlowFrom(OperatorFlow(mapOp(f), append(attributes, attribute.Name("MapFlow"))...))
}

func MapConcat(f MapConcatFunc, attributes ...attribute.StageAttribute) *stream.FlowGraph {
	return stream.FlowFrom(OperatorFlow(mapConcat(f), append(attributes, attribute.Name("MapFlow"))...))
}

func Filter(f FilterFunc, attributes ...attribute.StageAttribute) *stream.FlowGraph {
	return stream.FlowFrom(OperatorFlow(filter(f), append(attributes, attribute.Name("FilterFlow"))...))
}

func Fold(zero interface{}, f AggregateFunc, attributes ...attribute.StageAttribute) *stream.FlowGraph {
	return stream.FlowFrom(OperatorFlow(fold(zero, f), append(attributes, attribute.Name("FoldFlow"))...))
}

func Scan(zero interface{}, f AggregateFunc, attributes ...attribute.StageAttribute) *stream.FlowGraph {
	return stream.FlowFrom(OperatorFlow(scan(zero, f), append(attributes, attribute.Name("ScanFlow"))...))
}

func Take(number int, attributes ...attribute.StageAttribute) *stream.FlowGraph {
	return stream.FlowFrom(OperatorFlow(take(number), append(attributes, attribute.Name("TakeFlow"))...))
}

func TakeWhile(f FilterFunc, attributes ...attribute.StageAttribute) *stream.FlowGraph {
	return stream.FlowFrom(OperatorFlow(takeWhile(f), append(attributes, attribute.Name("TakeWhileFlow"))...))
}

func Concat(graphs ...*stream.FlowGraph) stream.FlowGraphFactory {
	return func(attributes ...attribute.StageAttribute) *stream.FlowGraph {
		return stream.CombineFlows(graphs)(stream.ConcatStrategy(), append(attributes, attribute.Name("ConcatFlow"))...)
	}
}

func Merge(graphs ...*stream.FlowGraph) stream.FlowGraphFactory {
	return func(attributes ...attribute.StageAttribute) *stream.FlowGraph {
		return stream.CombineFlows(graphs)(stream.MergeStrategy(), append(attributes, attribute.Name("MergeFlow"))...)
	}
}

func Interleave(segmentSize int, graphs ...*stream.FlowGraph) stream.FlowGraphFactory {
	return func(attributes ...attribute.StageAttribute) *stream.FlowGraph {
		return stream.CombineFlows(graphs)(stream.InterleaveStrategy(segmentSize), append(attributes, attribute.Name("InterleaveFlow"))...)
	}
}
