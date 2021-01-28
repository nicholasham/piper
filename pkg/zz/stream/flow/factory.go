package flow

import (
	"github.com/nicholasham/piper/pkg/zz/stream"
)

func Log(name string, attributes ...stream.StageOption) *stream.FlowGraph {
	return stream.FlowFrom(logFlow(name, attributes...))
}

func Operator(name string, operator OperatorLogic, attributes ...stream.StageOption) *stream.FlowGraph {
	return stream.FlowFrom(OperatorFlow(name, operator, attributes...))
}

func Map(f MapFunc, attributes ...stream.StageOption) *stream.FlowGraph {
	return Operator("MapFlow", mapOp(f), attributes...)
}

func MapConcat(f MapConcatFunc, attributes ...stream.StageOption) *stream.FlowGraph {
	return Operator("MapFlow", mapConcat(f), attributes...)
}

func Filter(f FilterFunc, attributes ...stream.StageOption) *stream.FlowGraph {
	return Operator("FilterFlow", filter(f), attributes...)
}

func Fold(zero interface{}, f AggregateFunc, attributes ...stream.StageOption) *stream.FlowGraph {
	return Operator("FoldFlow", fold(zero, f), attributes...)
}

func Scan(zero interface{}, f AggregateFunc, attributes ...stream.StageOption) *stream.FlowGraph {
	return Operator("ScanFlow", scan(zero, f), attributes...)
}

func Take(number int, attributes ...stream.StageOption) *stream.FlowGraph {
	return Operator("TakeFlow", take(number), attributes...)
}

func TakeWhile(f FilterFunc, attributes ...stream.StageOption) *stream.FlowGraph {
	return Operator("TakeWhileFlow", takeWhile(f), attributes...)
}

func Concat(graphs []*stream.FlowGraph, attributes ...stream.StageOption) *stream.FlowGraph {
	return stream.CombineFlows("ConcatFlow", graphs, stream.ConcatStrategy(), attributes...)
}

func Merge(graphs []*stream.FlowGraph, attributes ...stream.StageOption) *stream.FlowGraph {
	return stream.CombineFlows("MergeFlow", graphs, stream.MergeStrategy(), attributes...)
}

func Interleave(segmentSize int, graphs []*stream.FlowGraph, attributes ...stream.StageOption) *stream.FlowGraph {
	return stream.CombineFlows("InterleaveFlow", graphs, stream.InterleaveStrategy(segmentSize), attributes...)
}
