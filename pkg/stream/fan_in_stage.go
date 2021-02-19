package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
	"sync"
)

// verify fanInFlowStage implements FlowStage interface
var _ FlowStage = (*fanInFlowStage)(nil)

type FanInStrategy func(ctx context.Context, receivers []*Receiver, sender *Sender)

type fanInFlowStage struct {
	attributes     *StageAttributes
	upstreamStages []UpstreamStage
	fanIn          FanInStrategy
}

func (receiver *fanInFlowStage) Named(name string) Stage {
	return receiver.With(Name(name))
}

func (receiver *fanInFlowStage) Open(ctx context.Context, wg *sync.WaitGroup, mat MaterializeFunc) (*Receiver, *core.Future) {
	outputStream := NewStream(receiver.attributes.Name)
	outputPromise := core.NewPromise()

	var inlets []*Receiver

	for _, stage := range receiver.upstreamStages {
		reader, _ := stage.Open(ctx, wg, mat)
		inlets = append(inlets, reader)
	}

	receiver.fanIn(ctx, inlets, outputStream.Sender())

	outputPromise.TrySuccess(NotUsed)
	return outputStream.Receiver(), outputPromise.Future()
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
	return func(ctx context.Context, inletReaders []*Receiver, outletWriter *Sender) {
		go func() {
			defer outletWriter.Close()
			for _, inlet := range inletReaders {
				for element := range inlet.Receive() {

					select {
					case <-ctx.Done():
						inlet.Done()
					case <-outletWriter.Done():
						inlet.Done()
					default:

					}

					outletWriter.TrySend(element)

				}
			}
		}()
	}
}

func MergeStrategy() FanInStrategy {
	return func(ctx context.Context, receivers []*Receiver, senders *Sender) {
		wg := sync.WaitGroup{}
		wg.Add(len(receivers))

		f := func(receiver *Receiver, sender *Sender) {
			defer wg.Done()
			for element := range receiver.Receive() {

				select {
				case <-ctx.Done():
					receiver.Done()
				case <-sender.Done():
					receiver.Done()
				default:
				}

				sender.TrySend(element)

			}
		}

		go func() {
			for _, inlet := range receivers {
				go f(inlet, senders)
			}
		}()

		go func() {
			wg.Wait()
			senders.Close()
		}()
	}
}

func InterleaveStrategy(segmentSize int) FanInStrategy {
	return func(ctx context.Context, receivers []*Receiver, sender *Sender) {
		go func() {
			defer sender.Close()
			interleaveRecursively(ctx, segmentSize, receivers, sender)
		}()
	}
}

func interleaveRecursively(ctx context.Context, segmentSize int, receivers []*Receiver, sender *Sender) {
	var openReceivers []*Receiver
	for _, receiver := range receivers {
		if sendOutSegment(ctx, segmentSize, receiver, sender) {
			openReceivers = append(openReceivers, receiver)
		}
	}

	if len(openReceivers) > 0 {
		interleaveRecursively(ctx, segmentSize, openReceivers, sender)
	}
}

func sendOutSegment(ctx context.Context, segmentSize int, receiver *Receiver, sender *Sender) bool {
	segmentCount := 0
	for {
		select {
		case <-ctx.Done():
			return false
		case <-sender.Done():
			return false
		case element, ok := <-receiver.Receive():
			if !ok {
				return false
			}

			sender.TrySend(element)

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
