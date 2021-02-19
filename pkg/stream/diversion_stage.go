package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
	"sync"
)

type diversionStrategy func( element Element, mainWriter *Sender, alternateWriter *Sender)

var divertToStrategy = func(when core.PredicateFunc) diversionStrategy {
	return func(element Element, mainWriter *Sender, alternateWriter *Sender) {
		if element.IsValue() && when(element.Value()) {
			mainWriter.TrySend(element)
		} else {
			alternateWriter.TrySend(element)
		}
	}
}

var alsoToStrategy = func() diversionStrategy {
	return func(element Element, mainWriter *Sender, alternateWriter *Sender) {
		mainWriter.TrySend(element)
		alternateWriter.TrySend(element)
	}
}

// verify diversionFlowStage implements FlowStage interface
var _ FlowStage = (*diversionFlowStage)(nil)

type diversionFlowStage struct {
	attributes      *StageAttributes
	upstreamStage   UpstreamStage
	diversionSink   SinkStage
	diversionSource *diversionSourceStage
	strategy        diversionStrategy
}

func (d *diversionFlowStage) Named(name string) Stage {
	return d.With(Name(name))
}

func (d *diversionFlowStage) With(options ...StageOption) Stage {
	return &diversionFlowStage{
		attributes:      d.attributes.With(options...),
		upstreamStage:   d.upstreamStage,
		diversionSink:   d.diversionSink,
		diversionSource: d.diversionSource,
		strategy:        d.strategy,
	}
}

func (d *diversionFlowStage) Open(ctx context.Context, wg *sync.WaitGroup, mat MaterializeFunc) (*Receiver, *core.Future) {
	outputStream := NewStream(d.attributes.Name)
	outputPromise := core.NewPromise()
	reader, inputFuture := d.upstreamStage.Open(ctx, wg, KeepRight)
	wg.Add(1)
	go func() {
		diversionWriter := d.diversionSource.OpenWriter()
		mainWriter := outputStream.Sender()
		defer func() {
			diversionWriter.Close()
			mainWriter.Close()
			wg.Done()
		}()
		for element := range reader.Receive() {
			select {
			case <-ctx.Done():
				outputPromise.TryFailure(ctx.Err())
				reader.Done()
			case <-mainWriter.Done():
				reader.Done()
			default:
			}

			d.strategy(element, mainWriter, diversionWriter)
		}
	}()
	return outputStream.Receiver(), mat(inputFuture, outputPromise.Future())
}

func (d *diversionFlowStage) WireTo(stage UpstreamStage) FlowStage {
	d.upstreamStage = stage
	d.diversionSink.WireTo(d.diversionSource)
	return d
}

// verify diversionSourceStage implements UpstreamStage interface
var _ UpstreamStage = (*diversionSourceStage)(nil)

type diversionSourceStage struct {
	stream *Stream
}

func (d *diversionSourceStage) Open(ctx context.Context, wg *sync.WaitGroup, mat MaterializeFunc) (*Receiver, *core.Future) {
	promise := core.NewPromise()
	promise.TrySuccess(NotUsed)
	return d.stream.Receiver(), promise.Future()
}

func (d *diversionSourceStage) OpenWriter() *Sender {
	return d.stream.Sender()
}

func newDiversionStage(attributes  *StageAttributes) *diversionSourceStage {
	return & diversionSourceStage{stream: NewStream(attributes.Name + "-diversion")}
}


func diversion(diversionSink SinkStage, strategy diversionStrategy) FlowStage {
	attributes := DefaultStageAttributes.With(Name("diverted-stage"))
	return &diversionFlowStage{
		attributes:      attributes,
		upstreamStage:   nil,
		diversionSink:   diversionSink,
		diversionSource: newDiversionStage(attributes) ,
		strategy:        strategy,
	}
}