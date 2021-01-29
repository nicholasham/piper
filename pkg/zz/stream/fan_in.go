package stream

import (
	"context"
	"sync"
)

// verify fanInFlowStage implements stream.FlowStage interface
var _ FlowStage = (*fanInFlowStage)(nil)

type FanInStrategy func(ctx context.Context, inlets []*Inlet, outlet *Outlet)

type fanInFlowStage struct {
	name   string
	inlets []*Inlet
	outlet *Outlet
	fanIn  FanInStrategy
}

func (receiver *fanInFlowStage) Name() string {
	return receiver.name
}

func (receiver *fanInFlowStage) Run(ctx context.Context) {
	receiver.fanIn(ctx, receiver.inlets, receiver.outlet)
}

func (receiver *fanInFlowStage) Outlet() *Outlet {
	return receiver.outlet
}

func (receiver *fanInFlowStage) Wire(stage SourceStage) {
	inlet := NewInletOld(stage.Name()).WireTo(stage.Outlet())
	receiver.inlets = append(receiver.inlets, inlet)
}

func FanInFlow(name string, stages []SourceStage, strategy FanInStrategy, options ...StageOption) *fanInFlowStage {
	stageOptions := DefaultStageOptions.
		Apply(Name("AlsoToFlow")).
		Apply(options...)
	flow := fanInFlowStage{
		name:   stageOptions.Name,
		outlet: NewOutlet(stageOptions),
		fanIn:  strategy,
	}

	for _, stage := range stages {
		flow.Wire(stage)
	}

	return &flow
}

func CombineSources(name string, graphs []*SourceGraph, strategy FanInStrategy, options ...StageOption) *SourceGraph {
	var stages []SourceStage
	for _, graph := range graphs {
		stages = append(stages, graph.stage)
	}
	return SourceFrom(FanInFlow(name, stages, strategy, options...))
}

func CombineFlows(name string, graphs []*FlowGraph, strategy FanInStrategy, options ...StageOption) *FlowGraph {
	var stages []SourceStage
	for _, graph := range graphs {
		stages = append(stages, graph.stage)
	}
	return FlowFrom(FanInFlow(name, stages, strategy, options...))
}

func ConcatStrategy() FanInStrategy {
	return func(ctx context.Context, inlets []*Inlet, outlet *Outlet) {
		go func() {
			defer outlet.Close()
			for _, inlet := range inlets {
				for element := range inlet.In() {

					select {
					case <-ctx.Done():
						inlet.Complete()
					case <-outlet.Done():
						inlet.Complete()
					default:

					}

					if element.IsError() {
						outlet.Send(element)
						inlet.Complete()
					}

					if !inlet.CompletionSignaled() {
						outlet.Send(element)
					}
				}
			}
		}()
	}
}

func MergeStrategy() FanInStrategy {
	return func(ctx context.Context, inlets []*Inlet, outlet *Outlet) {
		wg := sync.WaitGroup{}
		wg.Add(len(inlets))

		f := func(inlet *Inlet, outlet *Outlet) {
			defer wg.Done()
			for element := range inlet.In() {

				select {
				case <-ctx.Done():
					inlet.Complete()
				case <-outlet.Done():
					inlet.Complete()
				default:
				}

				if element.IsError() {
					outlet.Send(element)
					inlet.Complete()
				}

				if !inlet.CompletionSignaled() {
					outlet.Send(element)
				}
			}
		}

		go func() {
			for _, inlet := range inlets {
				go f(inlet, outlet)
			}
		}()

		go func() {
			wg.Wait()
			outlet.Close()
		}()
	}
}

func InterleaveStrategy(segmentSize int) FanInStrategy {
	return func(ctx context.Context, inlets []*Inlet, outlet *Outlet) {
		go func() {
			defer outlet.Close()
			interleaveRecursively(ctx, segmentSize, inlets, outlet)
		}()
	}
}

func interleaveRecursively(ctx context.Context, segmentSize int, inlets []*Inlet, outlet *Outlet) {
	var openInlets []*Inlet
	for _, inlet := range inlets {
		if sendOutSegment(ctx, segmentSize, inlet, outlet) {
			openInlets = append(openInlets, inlet)
		}
	}

	if len(openInlets) > 0 {
		interleaveRecursively(ctx, segmentSize, openInlets, outlet)
	}
}

func sendOutSegment(ctx context.Context, segmentSize int, inlet *Inlet, outlet *Outlet) bool {
	segmentCount := 0
	for {
		select {
		case <-ctx.Done():
			return false
		case <-outlet.Done():
			return false
		case element, ok := <-inlet.in:
			if !ok {
				return false
			}

			outlet.Send(element)
			segmentCount++

			if segmentCount == segmentSize {
				return true
			}
		default:
		}
	}
}
