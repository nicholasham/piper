package stream

import "github.com/nicholasham/piper/pkg/stream/attribute"

type FlowGraphFactory func(attributes ...attribute.StageAttribute) *FlowGraph

type FlowGraph struct {
	stages []Stage
	stage  FlowStage
}

func (receiver *FlowGraph) Via(that *FlowGraph) *FlowGraph {
	that.stage.Wire(that.stage)
	that.stage.Wire(receiver.stage)
	return FlowFrom(that.stage, receiver.combineStages(that.stages)...)
}

func (receiver *FlowGraph) To(that *SinkGraph) *RunnableGraph {
	return flowRunnable(receiver, that)
}

func (receiver *FlowGraph) DivertTo(that *SinkGraph, predicate PredicateFunc, attributes ...attribute.StageAttribute) *FlowGraph {
	diversionStage := diversion(receiver.stage, that.stage, predicate, attributes)
	combinedStages := receiver.combineStages(that.stages)
	return FlowFrom(diversionStage, combinedStages...)
}

func (receiver *FlowGraph) AlsoTo(that *SinkGraph, attributes ...attribute.StageAttribute) *FlowGraph {
	diversionStage := alsoTo(receiver.stage, that.stage, attributes)
	combinedStages := receiver.combineStages(that.stages)
	return FlowFrom(diversionStage, combinedStages...)
}

func (receiver *FlowGraph) combineStages(stages []Stage) []Stage {
	var result []Stage
	for _, stage := range receiver.stages {
		result = append(result, stage)
	}
	for _, stage := range stages {
		result = append(result, stage)
	}
	return removeDuplicates(result)
}

func FlowFrom(stage FlowStage, stages ...Stage) *FlowGraph {
	return &FlowGraph{
		stages: append(stages, stage),
		stage:  stage,
	}
}
