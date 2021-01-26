package stream

import "context"

type FlowGraph struct {
	stage  FlowStage
}

func (receiver *FlowGraph) Via(that *FlowGraph) *FlowGraph {
	that.stage.Wire(receiver.stage)
	return FlowFrom(NewFusedFlow(receiver.stage, that.stage))
}

func (receiver *FlowGraph) To(that *SinkGraph) *RunnableGraph {
	return runnable(receiver.stage, that.stage)
}

func (receiver *FlowGraph) DivertTo(that *SinkGraph, predicate PredicateFunc, attributes ...StageAttribute) *FlowGraph {
	diversionStage := diversion(receiver.stage, that.stage, predicate, attributes)
	return FlowFrom(NewFusedFlow(receiver.stage, diversionStage))
}

func (receiver *FlowGraph) AlsoTo(that *SinkGraph, attributes ...StageAttribute) *FlowGraph {
	diversionStage := alsoTo(receiver.stage, that.stage, attributes)
	return FlowFrom(NewFusedFlow(receiver.stage, diversionStage))
}

func (receiver *FlowGraph) RunWith(ctx context.Context, that *SinkGraph) Future {
	return receiver.To(that).Run(ctx)
}

func FlowFrom(stage FlowStage) *FlowGraph {
	return &FlowGraph{
		stage:  stage,
	}
}
