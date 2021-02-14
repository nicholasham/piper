package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
)

type diversionStrategy func( element Element, mainWriter Writer, alternateWriter Writer)

var divertToStrategy = func(when core.PredicateFunc) diversionStrategy {
	return func(element Element, mainWriter Writer, alternateWriter Writer) {
		if element.IsValue() && when(element.Value()) {
			mainWriter.Send(element)
		} else {
			alternateWriter.Send(element)
		}
	}
}

var alsoToStrategy = func() diversionStrategy {
	return func(element Element, mainWriter Writer, alternateWriter Writer) {
		mainWriter.Send(element)
		alternateWriter.Send(element)
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

func (d *diversionFlowStage) Open(ctx context.Context, mat MaterializeFunc) (Reader, *core.Future) {
	outputStream := NewStream()
	outputPromise := core.NewPromise()
	reader, inputFuture := d.upstreamStage.Open(ctx, KeepRight)
	go func() {
		diversionWriter := d.diversionSource.OpenWriter()
		mainWriter := outputStream.Writer()
		for element := range reader.Read() {
			select {
			case <-ctx.Done():
				outputPromise.TryFailure(ctx.Err())
				reader.Complete()
				break
			case <-mainWriter.Done():
				reader.Complete()
				break
			default:
			}

			if !reader.Completing() {
				d.strategy(element, mainWriter, diversionWriter)
			}
		}
	}()
	return outputStream.Reader(), mat(inputFuture, outputPromise.Future())
}

func (d *diversionFlowStage) WireTo(stage UpstreamStage) FlowStage {
	d.upstreamStage = stage
	d.diversionSink.WireTo(d.diversionSource)
	return d
}

// verify diversionSourceStage implements UpstreamStage interface
var _ UpstreamStage = (*diversionSourceStage)(nil)

type diversionSourceStage struct {
	stream *Stream2
}

func (d *diversionSourceStage) Open(ctx context.Context, mat MaterializeFunc) (Reader, *core.Future) {
	promise := core.NewPromise()
	promise.TrySuccess(NotUsed)
	return d.stream.Reader(), promise.Future()
}

func (d *diversionSourceStage) OpenWriter() Writer {
	return d.stream.Writer()
}

func newDiversionStage(attributes  *StageAttributes) *diversionSourceStage {
	return & diversionSourceStage{stream: NewStream()}
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