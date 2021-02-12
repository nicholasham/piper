package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
	"sync"
)

// verify fanInFlowStage implements FlowStage interface
var _ FlowStage = (*fanInFlowStage)(nil)

type FanInStrategy func(ctx context.Context, inletReaders []Reader, outletWriter Writer)

type fanInFlowStage struct {
	attributes     *StageAttributes
	upstreamStages []UpstreamStage
	fanIn          FanInStrategy
}

func (receiver *fanInFlowStage) Named(name string) Stage {
	return receiver.With(Name(name))
}

func (receiver *fanInFlowStage) Open(ctx context.Context, mat MaterializeFunc) (Reader, *core.Future) {
	outputStream := NewStream(receiver.attributes.Name)
	outputPromise := core.NewPromise()

	var inlets []Reader

	for _, stage := range receiver.upstreamStages {
		reader, _ := stage.Open(ctx, mat)
		inlets = append(inlets, reader)
	}

	receiver.fanIn(ctx, inlets, outputStream.Writer())

	outputPromise.TrySuccess(NotUsed)
	return outputStream.Reader(), outputPromise.Future()
}

func (receiver *fanInFlowStage) WireTo(stage UpstreamStage) FlowStage {
	receiver.upstreamStages = append(receiver.upstreamStages, stage)
	return receiver
}

func (receiver *fanInFlowStage) With(opts ...StageOption) Stage {
	options := receiver.attributes.With(opts...)
	return &fanInFlowStage{
		attributes:     options,
		upstreamStages: receiver.upstreamStages,
		fanIn:          receiver.fanIn,
	}
}

func FanInFlow(stages []SourceStage, strategy FanInStrategy) FlowStage {
	attributes := DefaultStageAttributes.With(Name("FanIn"))
	flow := fanInFlowStage{
		attributes:     attributes,
		fanIn:          strategy,
		upstreamStages: sourcesToUpstream(stages),
	}
	return &flow
}

func sourcesToUpstream(stages []SourceStage) []UpstreamStage {
	var upstreamStages []UpstreamStage
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
	return func(ctx context.Context, inletReaders []Reader, outletWriter Writer) {
		go func() {
			defer outletWriter.Close()
			for _, inlet := range inletReaders {
				for element := range inlet.Elements() {

					select {
					case <-ctx.Done():
						inlet.Complete()
					case <-outletWriter.Done():
						inlet.Complete()
					default:

					}

					if !inlet.Completing() {
						element.WhenError(outletWriter.SendError)
						element.WhenValue(outletWriter.SendValue)
					}
				}
			}
		}()
	}
}

func MergeStrategy() FanInStrategy {
	return func(ctx context.Context, inletReaders []Reader, outletReader Writer) {
		wg := sync.WaitGroup{}
		wg.Add(len(inletReaders))

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
			for _, inlet := range inletReaders {
				go f(inlet, outletReader)
			}
		}()

		go func() {
			wg.Wait()
			outletReader.Close()
		}()
	}
}

func InterleaveStrategy(segmentSize int) FanInStrategy {
	return func(ctx context.Context, inletReaders []Reader, outletWriter Writer) {
		go func() {
			defer outletWriter.Close()
			interleaveRecursively(ctx, segmentSize, inletReaders, outletWriter)
		}()
	}
}

func interleaveRecursively(ctx context.Context, segmentSize int, inletReaders []Reader, outletWriter Writer) {
	var openInletReaders []Reader
	for _, inletReader := range inletReaders {
		if sendOutSegment(ctx, segmentSize, inletReader, outletWriter) {
			openInletReaders = append(openInletReaders, inletReader)
		}
	}

	if len(openInletReaders) > 0 {
		interleaveRecursively(ctx, segmentSize, openInletReaders, outletWriter)
	}
}

func sendOutSegment(ctx context.Context, segmentSize int, inletReader Reader, outletWriter Writer) bool {
	segmentCount := 0
	for {
		select {
		case <-ctx.Done():
			return false
		case <-outletWriter.Done():
			return false
		case element, ok := <-inletReader.Elements():
			if !ok {
				return false
			}

			if !inletReader.Completing() {
				element.WhenError(outletWriter.SendError)
				element.WhenValue(outletWriter.SendValue)
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
