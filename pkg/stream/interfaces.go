package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
)

type Stage interface {
	Name() string
	Run(ctx context.Context)
	With(options ...StageOption) Stage
}

type SinkStage interface {
	InputStage
	WireTo(stage OutputStage) SinkStage
	Result() Future
}

type Future interface {
	Await() core.Result
}

type OutputStage interface {
	Stage
	Outlet() *Outlet
}

type InputStage interface {
	Stage
}

type SourceStage interface {
	OutputStage
}

type FlowStage interface {
	InputStage
	OutputStage
	WireTo(stage OutputStage) FlowStage
}

type OutStageLogic interface {
	OnPull(actions OutStageActions)
	OnDownstreamFinish(actions OutStageActions)
}

type OutStageActions interface {
	Push(element Element)
}

type InStageLogic interface {
	// Called when starting to receive elements from upstream
	OnUpstreamStart(actions InStageActions)
	// Called when an element is received from upstream
	OnPush(value interface{}, actions InStageActions)
	// Called when up stream has failed
	OnUpstreamFailure(cause error, actions InStageActions)
	// 	Called when finishing receiving elements from upstream
	OnUpstreamFinish(actions InStageActions)
}

type InStageActions interface {
	// Sends an error downstream
	SendError(cause error)
	// Sends a value downstream
	SendValue(value interface{})
	// Fails a stage on logs the cause of failure.
	FailStage(cause error)
	// Completes the stage
	CompleteStage()
}
