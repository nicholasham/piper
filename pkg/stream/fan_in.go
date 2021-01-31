package stream

import (
	"context"
	"sync"
)

// verify fanInFlowStage implements FlowStage interface
var _ FlowStage = (*fanInFlowStage)(nil)

type FanInStrategy func(ctx context.Context, inlets []*Inlet, outlet *Outlet)

type fanInFlowStage struct {
	attributes *StageAttributes
	inlets     []*Inlet
	outlet     *Outlet
	fanIn      FanInStrategy
}

func (receiver *fanInFlowStage) WireTo(stage OutputStage) FlowStage {
	inlet := NewInlet(receiver.attributes)
	inlet.WireTo(stage.Outlet())
	receiver.inlets = append(receiver.inlets, inlet)
	return receiver
}

func (receiver *fanInFlowStage) With(opts ...StageOption) Stage {
	options := receiver.attributes.Apply(opts...)
	return &fanInFlowStage{
		attributes: options,
		inlets:     receiver.inlets,
		outlet:     NewOutlet(options),
		fanIn:      receiver.fanIn,
	}
}

func (receiver *fanInFlowStage) Name() string {
	return receiver.attributes.Name
}

func (receiver *fanInFlowStage) Run(ctx context.Context) {
	receiver.fanIn(ctx, receiver.inlets, receiver.outlet)
}

func (receiver *fanInFlowStage) Outlet() *Outlet {
	return receiver.outlet
}

func FanInFlow(stages []SourceStage, strategy FanInStrategy) FlowStage {
	attributes := DefaultStageAttributes.Apply(Name("FanIn"))
	flow := fanInFlowStage{
		attributes: attributes,
		outlet:     NewOutlet(attributes),
		fanIn:      strategy,
	}

	for _, stage := range stages {
		flow.WireTo(stage)
	}

	return &flow
}

func CombineSources(graphs []*SourceGraph, strategy FanInStrategy) *SourceGraph {
	var stages []SourceStage
	for _, graph := range graphs {
		stages = append(stages, graph.stage)
	}
	return FromSource(FanInFlow(stages, strategy))
}

func CombineFlows(graphs []*FlowGraph, strategy FanInStrategy) *FlowGraph {
	var stages []SourceStage
	for _, graph := range graphs {
		stages = append(stages, graph.stage)
	}
	return FromFlow(FanInFlow(stages, strategy))
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

					if !inlet.CompletionSignaled() {
						element.WhenError(outlet.SendError)
						element.WhenValue(outlet.SendValue)
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

				if !inlet.CompletionSignaled() {
					element.WhenError(outlet.SendError)
					element.WhenValue(outlet.SendValue)
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

			if !inlet.CompletionSignaled() {
				element.WhenError(outlet.SendError)
				element.WhenValue(outlet.SendValue)
			}

			segmentCount++

			if segmentCount == segmentSize {
				return true
			}
		default:
		}
	}
}

func ConcatSources(graphs ...*SourceGraph) *SourceGraph {
	return CombineSources(graphs, ConcatStrategy()).Named("ConcatSource")
}

func InterleaveSources(segmentSize int, graphs ...*SourceGraph) *SourceGraph {
	return CombineSources(graphs, InterleaveStrategy(segmentSize)).Named("InterleaveSource")
}

func MergeSources(graphs ...*SourceGraph) *SourceGraph {
	return CombineSources(graphs, MergeStrategy()).Named("MergeSource")
}

func ConcatFlows(graphs ...*FlowGraph) *FlowGraph {
	return CombineFlows(graphs, ConcatStrategy()).Named("ConcatFlows")
}

func InterleaveFlows(segmentSize int, graphs ...*FlowGraph) *FlowGraph {
	return CombineFlows(graphs, InterleaveStrategy(segmentSize)).Named("InterleaveSource")
}

func MergeFlows(graphs ...*FlowGraph) *FlowGraph {
	return CombineFlows(graphs, MergeStrategy()).Named("MergeSource")
}
