package piper

type FlowGraph struct {
	stages []Stage
	stage  FlowStage
}

func (receiver *FlowGraph) Via(that *FlowGraph) *FlowGraph {
	that.stage.Wire(receiver.stage)
	combinedStages := combineStages(receiver.stages, that.stages)
	return FlowFrom(that.stage, combinedStages...)
}

func (receiver *FlowGraph) To(that *SinkGraph) *RunnableGraph {
	return flowRunnable(receiver, that)
}

func (receiver *FlowGraph) DivertTo(that *SinkGraph, predicate PredicateFunc, attributes ...StageAttribute) *FlowGraph {
	diversionStage := diversion(receiver.stage, that.stage, predicate, attributes)
	combinedStages := combineStages(receiver.stages, that.stages)
	return FlowFrom(diversionStage, combinedStages...)
}

func (receiver *FlowGraph) AlsoTo(that *SinkGraph, attributes ...StageAttribute) *FlowGraph {
	diversionStage := alsoTo(receiver.stage, that.stage, attributes)
	combinedStages := combineStages(receiver.stages, that.stages)
	return FlowFrom(diversionStage, combinedStages...)
}


func FlowFrom(stage FlowStage, stages ...Stage) *FlowGraph {
	return &FlowGraph{
		stages: append(stages, stage),
		stage:  stage,
	}
}
