package flow

import (
	"github.com/nicholasham/piper/pkg/piper"
)

func Log(name string, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(logFlow(name, attributes...))
}

func Operator(name string, operator OperatorLogic, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(name, operator, attributes...))
}

func Map(f MapFunc, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return Operator("MapFlow", mapOp(f), attributes...)
}

func MapConcat(f MapConcatFunc, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return Operator("MapFlow", mapConcat(f), attributes...)
}

func Filter(f FilterFunc, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return Operator("FilterFlow", filter(f), attributes...)
}

func Fold(zero interface{}, f AggregateFunc, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return Operator("FoldFlow", fold(zero, f), attributes...)
}

func Scan(zero interface{}, f AggregateFunc, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return Operator("ScanFlow", scan(zero, f), attributes...)
}

func Take(number int, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return Operator("TakeFlow", take(number), attributes...)
}

func TakeWhile(f FilterFunc, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return Operator("TakeWhileFlow", takeWhile(f), attributes...)
}

func Concat(graphs []*piper.FlowGraph, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.CombineFlows("ConcatFlow", graphs, piper.ConcatStrategy(), attributes...)
}

func Merge(graphs []*piper.FlowGraph, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.CombineFlows("MergeFlow", graphs, piper.MergeStrategy(), attributes...)
}

func Interleave(segmentSize int, graphs []*piper.FlowGraph, attributes ...piper.StageAttribute) *piper.FlowGraph {
	return piper.CombineFlows("InterleaveFlow", graphs, piper.InterleaveStrategy(segmentSize), attributes...)
}
