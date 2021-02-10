package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
	"sync"
)

// verify fanInFlowStage implements FlowStage interface
var _ FlowStage = (*fanInFlowStage)(nil)

type FanInStrategy func(ctx context.Context, readers []Reader, writers Writer)

type fanInFlowStage struct {
	attributes *StageAttributes
	upstreamStages[]UpstreamStage
	fanIn      FanInStrategy
}

func (receiver *fanInFlowStage) Open(ctx context.Context, mat MaterializeFunc) (Reader, *core.Future) {
	panic("implement me")
}

func (receiver *fanInFlowStage) WireTo(stage UpstreamStage) FlowStage {
	receiver.upstreamStages = append(receiver.upstreamStages, stage)
	return receiver
}

func (receiver *fanInFlowStage) With(opts ...StageOption) Stage {
	options := receiver.attributes.Apply(opts...)
	return &fanInFlowStage{
		attributes: options,
		upstreamStages:     receiver.upstreamStages,
		fanIn:      receiver.fanIn,
	}
}

func (receiver *fanInFlowStage) Name() string {
	return receiver.attributes.Name
}

func (receiver *fanInFlowStage) Run(ctx context.Context) {

	//receiver.fanIn(ctx, receiver.inlets, receiver.outlet)
}

func FanInFlow(stages []SourceStage, strategy FanInStrategy) FlowStage {
	attributes := DefaultStageAttributes.Apply(Name("FanIn"))
	flow := fanInFlowStage{
		attributes: attributes,
		fanIn:      strategy,
		upstreamStages: sourcesToUpstream(stages),
	}
	return &flow
}

func sourcesToUpstream(stages []SourceStage) []UpstreamStage {
	var upstreamStages[] UpstreamStage
	for _, stage := range stages {
		upstreamStages = append(upstreamStages, stage)
	}
	return upstreamStages
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
	return func(ctx context.Context, inlets []Reader, outlet Writer) {
		go func() {
			defer outlet.Close()
			for _, inlet := range inlets {
				for element := range inlet.Elements() {

					select {
					case <-ctx.Done():
						inlet.Complete()
					case <-outlet.Done():
						inlet.Complete()
					default:

					}

					if !inlet.Completing() {
						element.WhenError(outlet.SendError)
						element.WhenValue(outlet.SendValue)
					}
				}
			}
		}()
	}
}

func MergeStrategy() FanInStrategy {
	return func(ctx context.Context, inlets []Reader, outlet Writer) {
		wg := sync.WaitGroup{}
		wg.Add(len(inlets))

		f := func(inlet Reader, outlet Writer) {
			defer wg.Done()
			for element := range inlet.Elements() {

				select {
				case <-ctx.Done():
					inlet.Complete()
				case <-outlet.Done():
					inlet.Complete()
				default:
				}

				if !inlet.Completing() {
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
	return func(ctx context.Context, inlets []Reader, outlet Writer) {
		go func() {
			defer outlet.Close()
			interleaveRecursively(ctx, segmentSize, inlets, outlet)
		}()
	}
}

func interleaveRecursively(ctx context.Context, segmentSize int, inlets []Reader, outlet Writer) {
	var openInlets []Reader
	for _, inlet := range inlets {
		if sendOutSegment(ctx, segmentSize, inlet, outlet) {
			openInlets = append(openInlets, inlet)
		}
	}

	if len(openInlets) > 0 {
		interleaveRecursively(ctx, segmentSize, openInlets, outlet)
	}
}

func sendOutSegment(ctx context.Context, segmentSize int, inlet Reader, outlet Writer) bool {
	segmentCount := 0
	for {
		select {
		case <-ctx.Done():
			return false
		case <-outlet.Done():
			return false
		case element, ok := <-inlet.Elements():
			if !ok {
				return false
			}

			if !inlet.Completing() {
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
