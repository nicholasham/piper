package piper

import (
	"context"
)

type SourceGraph struct {
	stage  SourceStage
	stages []Stage
}

func (receiver *SourceGraph) RunWith(ctx context.Context, that *SinkGraph) Future {
	return receiver.To(that).Run(ctx)
}

func (receiver *SourceGraph) DivertTo(that *SinkGraph, predicate PredicateFunc, attributes ...StageAttribute) *SourceGraph {
	diversionStage := diversion(receiver.stage, that.stage, predicate, attributes)
	combinedStages := combineStages(receiver.stages, that.stages)
	return SourceFrom(diversionStage, combinedStages...)
}

func (receiver *SourceGraph) AlsoTo(that *SinkGraph, attributes ...StageAttribute) *SourceGraph {
	diversionStage := alsoTo(receiver.stage, that.stage, attributes)
	combinedStages := combineStages(receiver.stages, that.stages)
	return SourceFrom(diversionStage, combinedStages...)
}

// Transform this FlowStage by appending the given processing steps.
func (receiver *SourceGraph) Via(that *FlowGraph) *FlowGraph {
	that.stage.Wire(receiver.stage)
	combinedStages := combineStages(receiver.stages, that.stages)
	return FlowFrom(that.stage, combinedStages...)
}

func (receiver *SourceGraph) To(that *SinkGraph) *RunnableGraph {
	return sourceRunnable(receiver, that)
}

func SourceFrom(sourceStage SourceStage, stages ...Stage) *SourceGraph {
	return &SourceGraph{
		stage:  sourceStage,
		stages: removeDuplicates(append(stages, sourceStage))}
}
