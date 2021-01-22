package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/stream/attribute"
)

type SourceGraphFactory func(attributes ...attribute.StageAttribute) *SourceGraph

type SourceGraph struct {
	stage  SourceStage
	stages []Stage
}

func (receiver *SourceGraph) RunWith(ctx context.Context, that *SinkGraph) Future {
	return receiver.To(that).Run(ctx)
}

func (receiver *SourceGraph) DivertTo(that *SinkGraph, predicate PredicateFunc, attributes ...attribute.StageAttribute) *SourceGraph {
	diversionStage := diversion(receiver.stage, that.stage, predicate, attributes)
	combinedStages := receiver.combineStages(that.stages)
	return SourceFrom(diversionStage, combinedStages...)
}

func (receiver *SourceGraph) AlsoTo(that *SinkGraph, attributes ...attribute.StageAttribute) *SourceGraph {
	diversionStage := alsoTo(receiver.stage, that.stage, attributes)
	combinedStages := receiver.combineStages(that.stages)
	return SourceFrom(diversionStage, combinedStages...)
}

func (receiver *SourceGraph) combineStages(stages []Stage) []Stage {
	var result []Stage
	for _, stage := range receiver.stages {
		result = append(result, stage)
	}
	for _, stage := range stages {
		result = append(result, stage)
	}
	return removeDuplicates(result)
}

// Transform this FlowStage by appending the given processing steps.
func (receiver *SourceGraph) Via(that *FlowGraph) *FlowGraph {
	that.stage.Wire(receiver.stage)
	return FlowFrom(that.stage, receiver.combineStages(that.stages)...)
}

func (receiver *SourceGraph) To(that *SinkGraph) *RunnableGraph {
	return sourceRunnable(receiver, that)
}

func SourceFrom(sourceStage SourceStage, stages ...Stage) *SourceGraph {
	return &SourceGraph{
		stage:  sourceStage,
		stages: append(stages, sourceStage)}
}
