package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
)

type diversionStrategy func( element Element, mainWriter Writer, alternateWriter Writer)

var divertToStrategy = func(when core.PredicateFunc) diversionStrategy {
	return func(element Element, mainWriter Writer, alternateWriter Writer) {
		if element.IsValue() && when(element.Value()) {
			element.
				WhenValue(mainWriter.SendValue)
		} else {
			element.
				WhenValue(mainWriter.SendValue).
				WhenError(mainWriter.SendError)
		}
	}
}

var alsoToStrategy = func() diversionStrategy {
	return func(element Element, mainWriter Writer, alternateWriter Writer) {
		element.
			WhenValue(mainWriter.SendValue).
			WhenValue(alternateWriter.SendValue).
			WhenError(mainWriter.SendError).
			WhenError(alternateWriter.SendError)
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

func (d *diversionFlowStage) With(options ...StageOption) Stage {
	return &diversionFlowStage{
		attributes:      d.attributes.Apply(options...),
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
		defer diversionWriter.Close()
		defer mainWriter.Close()
		for element := range reader.Elements() {
			select {
			case <-ctx.Done():
				outputPromise.TryFailure(ctx.Err())
				reader.Complete()
			case <-mainWriter.Done():
				reader.Complete()
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
	stream Stream
}

func (d *diversionSourceStage) Open(ctx context.Context, mat MaterializeFunc) (Reader, *core.Future) {
	promise := core.NewPromise()
	promise.TrySuccess(NotUsed)
	return d.stream.Reader(), promise.Future()
}

func (d *diversionSourceStage) OpenWriter() Writer {
	return d.stream.Writer()
}

func newDiversionStage() *diversionSourceStage {
	return & diversionSourceStage{stream: NewStream()}
}


func diversion(diversionSink SinkStage, strategy diversionStrategy) FlowStage {
	return &diversionFlowStage{
		attributes:      DefaultStageAttributes,
		upstreamStage:   nil,
		diversionSink:   diversionSink,
		diversionSource: newDiversionStage(),
		strategy:        strategy,
	}
}