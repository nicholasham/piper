package stream

import (
	"context"
)

type SourceGraph struct {
	stage  SourceStage
}

func (receiver *SourceGraph) RunWith(ctx context.Context, that *SinkGraph) Future {
	return receiver.To(that).Run(ctx)
}

func (receiver *SourceGraph) DivertTo(that *SinkGraph, predicate PredicateFunc, attributes ...StageAttribute) *SourceGraph {
	diversionStage := diversion(receiver.stage, that.stage, predicate, attributes)
	return SourceFrom(NewCompositeFlow(receiver.stage, diversionStage))
}

func (receiver *SourceGraph) AlsoTo(that *SinkGraph, attributes ...StageAttribute) *SourceGraph {
	diversionStage := alsoTo(receiver.stage, that.stage, attributes)
	return SourceFrom(NewCompositeFlow(receiver.stage, diversionStage))
}

// Transform this FlowStage by appending the given processing steps.
func (receiver *SourceGraph) Via(that *FlowGraph) *FlowGraph {
	return FlowFrom(NewCompositeFlow(receiver.stage, that.stage))
}

func (receiver *SourceGraph) To(that *SinkGraph) *RunnableGraph {
	return runnable(receiver.stage, that.stage)
}

func SourceFrom(sourceStage SourceStage) *SourceGraph {
	return &SourceGraph{
		stage:  sourceStage,
	}
}
