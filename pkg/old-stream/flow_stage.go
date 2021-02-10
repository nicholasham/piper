package old_stream

import (
	"context"

	"github.com/gammazero/workerpool"
)

// verify linearFlowStage implements old-stream.FlowStage interface
var _ FlowStage = (*linearFlowStage)(nil)
var _ FlowStageActions = (*linearFlowStage)(nil)

type FlowStageLogicFactory func(attributes *StageAttributes) FlowStageLogic

type FlowStageLogic interface {
	SupportsParallelism() bool
	// Called when starting to receive elements from upstream
	OnUpstreamStart(actions FlowStageActions)
	// Called when an element is received from upstream
	OnUpstreamReceive(element Element, actions FlowStageActions)
	// 	Called when finishing receiving elements from upstream
	OnUpstreamFinish(actions FlowStageActions)
}

type FlowStageActions interface {
	// Sends an error downstream
	SendError(cause error)
	// Sends a value downstream
	SendValue(value interface{})
	// Fails a stage on logs the cause of failure.
	FailStage(cause error)
	// Completes the stage
	CompleteStage()
}

type linearFlowStage struct {
	attributes *StageAttributes
	inlet      *Inlet
	outlet     *Outlet
	factory    FlowStageLogicFactory
}

func (o *linearFlowStage) SendError(cause error) {
	o.outlet.SendError(cause)
}

func (o *linearFlowStage) SendValue(value interface{}) {
	o.outlet.SendValue(value)
}

func (o *linearFlowStage) FailStage(cause error) {
	o.attributes.Logger.Error(cause, "failed stage because")
	o.inlet.Complete()
}

func (o *linearFlowStage) CompleteStage() {
	o.inlet.Complete()
}

func (o *linearFlowStage) With(options ...StageOption) Stage {
	attributes := o.attributes.Apply(options...)
	return &linearFlowStage{
		attributes: attributes,
		inlet:      NewInlet(attributes),
		outlet:     NewOutlet(attributes),
		factory:    o.factory,
	}
}

func (o *linearFlowStage) WireTo(stage OutputStage) FlowStage {
	o.inlet.WireTo(stage.Outlet())
	return o
}

func (o *linearFlowStage) Name() string {
	return o.attributes.Name
}

func (o *linearFlowStage) Run(ctx context.Context) {
	go func() {
		logic := o.factory(o.attributes)
		wp := o.createWorkerPool(logic)
		defer func() {
			o.outlet.Close()
		}()
		logic.OnUpstreamStart(o)
		for element := range o.inlet.In() {

			select {
			case <-ctx.Done():
				o.outlet.SendError(ctx.Err())
				o.inlet.Complete()
			case <-o.outlet.Done():
				o.inlet.Complete()
			default:
			}

			if !o.inlet.CompletionSignaled() {
				wp.Submit(o.Push(logic, element))
			}

		}
		wp.StopWait()
		logic.OnUpstreamFinish(o)
	}()
}

func (o *linearFlowStage) createWorkerPool(logic FlowStageLogic) *workerpool.WorkerPool {
	maxWorkers := 1
	if logic.SupportsParallelism() {
		maxWorkers = o.attributes.Parallelism
	}
	return workerpool.New(maxWorkers)
}

func (o *linearFlowStage) Push(logic FlowStageLogic, element Element) func() {
	return func() {
		logic.OnUpstreamReceive(element, o)
	}
}

func (o *linearFlowStage) Outlet() *Outlet {
	return o.outlet
}

func (o *linearFlowStage) Wire(stage SourceStage) {
	o.inlet.WireTo(stage.Outlet())
}

func LinearFlow(factory FlowStageLogicFactory) FlowStage {
	attributes := DefaultStageAttributes.Apply(Name("LinearFlowStage"))
	return &linearFlowStage{
		attributes: attributes,
		factory:    factory,
		inlet:      NewInlet(attributes),
		outlet:     NewOutlet(attributes),
	}
}
