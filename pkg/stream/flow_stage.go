package stream

import (
	"context"
	"github.com/gammazero/workerpool"
)

// verify flowStage implements stream.FlowStage interface
var _ FlowStage = (*flowStage)(nil)
var _ FlowStageActions = (*flowStage)(nil)

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
	SendError(cause error)
	SendValue(value interface{})
	FailStage(cause error)
	CompleteStage()
}

type flowStage struct {
	attributes *StageAttributes
	inlet      *Inlet
	outlet     *Outlet
	factory    FlowStageLogicFactory
}

func (o *flowStage) SendError(cause error) {
	o.outlet.SendError(cause)
}

func (o *flowStage) SendValue(value interface{}) {
	o.outlet.SendValue(value)
}

func (o *flowStage) FailStage(cause error) {
	o.attributes.Logger.Error(cause, "failed stage because")
	o.inlet.Complete()
}

func (o *flowStage) CompleteStage() {
	o.inlet.Complete()
}

func (o *flowStage) With(options ...StageOption) Stage {
	attributes := o.attributes.Apply(options...)
	return &flowStage{
		attributes: attributes,
		inlet:      NewInlet(attributes),
		outlet:     NewOutlet(attributes),
		factory:    o.factory,
	}
}

func (o *flowStage) WireTo(stage OutputStage) FlowStage {
	o.inlet.WireTo(stage.Outlet())
	return o
}

func (o *flowStage) Name() string {
	return o.attributes.Name
}

func (o *flowStage) Run(ctx context.Context) {
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

func (o *flowStage) createWorkerPool(logic FlowStageLogic) * workerpool.WorkerPool {
	maxWorkers := 1
	if logic.SupportsParallelism(){
		maxWorkers = o.attributes.Parallelism
	}
	return workerpool.New(maxWorkers)
}

func (o *flowStage) Push(logic FlowStageLogic, element Element) func() {
	return func() {
		logic.OnUpstreamReceive(element, o)
	}
}

func (o *flowStage) Outlet() *Outlet {
	return o.outlet
}

func (o *flowStage) Wire(stage SourceStage) {
	o.inlet.WireTo(stage.Outlet())
}

func LinearFlow(factory FlowStageLogicFactory) FlowStage {
	attributes := DefaultStageAttributes.Apply(Name("LinearFlowStage"))
	return &flowStage{
		attributes: attributes,
		factory:    factory,
		inlet:      NewInlet(attributes),
		outlet:     NewOutlet(attributes),
	}
}





