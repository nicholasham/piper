package flow

import (
	"github.com/nicholasham/piper/pkg/piper"
	"github.com/nicholasham/piper/pkg/piper/attribute"
)

func Log(name string, attributes ...attribute.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(logFlow(name, attributes...))
}

func Map(f MapFunc, attributes ...attribute.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(mapOp(f), append(attributes, attribute.Name("MapFlow"))...))
}

func MapConcat(f MapConcatFunc, attributes ...attribute.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(mapConcat(f), append(attributes, attribute.Name("MapFlow"))...))
}

func Filter(f FilterFunc, attributes ...attribute.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(filter(f), append(attributes, attribute.Name("FilterFlow"))...))
}

func Fold(zero interface{}, f AggregateFunc, attributes ...attribute.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(fold(zero, f), append(attributes, attribute.Name("FoldFlow"))...))
}

func Scan(zero interface{}, f AggregateFunc, attributes ...attribute.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(scan(zero, f), append(attributes, attribute.Name("ScanFlow"))...))
}

func Take(number int, attributes ...attribute.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(take(number), append(attributes, attribute.Name("TakeFlow"))...))
}

func TakeWhile(f FilterFunc, attributes ...attribute.StageAttribute) *piper.FlowGraph {
	return piper.FlowFrom(OperatorFlow(takeWhile(f), append(attributes, attribute.Name("TakeWhileFlow"))...))
}

func Concat(graphs ...*piper.FlowGraph) piper.FlowGraphFactory {
	return func(attributes ...attribute.StageAttribute) *piper.FlowGraph {
		return piper.CombineFlows(graphs)(piper.ConcatStrategy(), append(attributes, attribute.Name("ConcatFlow"))...)
	}
}

func Merge(graphs ...*piper.FlowGraph) piper.FlowGraphFactory {
	return func(attributes ...attribute.StageAttribute) *piper.FlowGraph {
		return piper.CombineFlows(graphs)(piper.MergeStrategy(), append(attributes, attribute.Name("MergeFlow"))...)
	}
}

func Interleave(segmentSize int, graphs ...*piper.FlowGraph) piper.FlowGraphFactory {
	return func(attributes ...attribute.StageAttribute) *piper.FlowGraph {
		return piper.CombineFlows(graphs)(piper.InterleaveStrategy(segmentSize), append(attributes, attribute.Name("InterleaveFlow"))...)
	}
}
