package flow

import (
	"github.com/nicholasham/piper/pkg/piper"
)

func Log(name string, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(logFlow(name, attributes...))
}

func Map(f MapFunc, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(mapOp(f), append(attributes, piper.Name("MapFlow"))...))
}

func MapConcat(f MapConcatFunc, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(mapConcat(f), append(attributes, piper.Name("MapFlow"))...))
}

func Filter(f FilterFunc, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(filter(f), append(attributes, piper.Name("FilterFlow"))...))
}

func Fold(zero interface{}, f AggregateFunc, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(fold(zero, f), append(attributes, piper.Name("FoldFlow"))...))
}

func Scan(zero interface{}, f AggregateFunc, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(scan(zero, f), append(attributes, piper.Name("ScanFlow"))...))
}

func Take(number int, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(take(number), append(attributes, piper.Name("TakeFlow"))...))
}

func TakeWhile(f FilterFunc, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(takeWhile(f), append(attributes, piper.Name("TakeWhileFlow"))...))
}

func Concat(graphs []*piper.FlowGraph, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.CombineFlows("ConcatFlow",graphs, piper.ConcatStrategy(), attributes...)
}

func Merge(graphs []*piper.FlowGraph, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.CombineFlows("MergeFlow", graphs, piper.MergeStrategy(), attributes...)
}

func Interleave(segmentSize int, graphs []*piper.FlowGraph, attributes ...piper.StageAttribute) *piper.FlowGraph  {
	return piper.CombineFlows("InterleaveFlow", graphs, piper.InterleaveStrategy(segmentSize), attributes...)
}
