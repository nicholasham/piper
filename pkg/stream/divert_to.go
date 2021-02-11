package stream

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
)

// verify divertToFlowStage implements FlowStage interface
var _ FlowStage = (*divertToFlowStage)(nil)

type divertToFlowStage struct {
	attributes *StageAttributes
	upstreamStage UpstreamStage
	diversionSink SinkStage
	divert core.PredicateFunc
}

func (d *divertToFlowStage) With(options ...StageOption) Stage {
	return &divertToFlowStage{
		attributes:    d.attributes.Apply(options...),
		upstreamStage: d.upstreamStage,
		diversionSink: d.diversionSink,
		divert:        d.divert,
	}
}

func (d *divertToFlowStage) Open(ctx context.Context, mat MaterializeFunc) (Reader, *core.Future) {
	outputStream := NewStream()
	outputPromise := core.NewPromise()
	reader, inputFuture := d.upstreamStage.Open(ctx, KeepRight)
	go func() {
		diversionWriter := d.createDiversionWriter()
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
				if element.IsValue() && d.divert(element.Value()){
					element.
						WhenValue(mainWriter.SendValue)
				}else {
					element.
						WhenValue(mainWriter.SendValue).
						WhenError(mainWriter.SendError)
				}
			}
		}
	}()
	return outputStream.Reader(), mat(inputFuture, outputPromise.Future())
}

func (d *divertToFlowStage) createDiversionWriter() Writer {
	diversionUpstream := &diversionUpstream{stream: NewStream()}
	diversionStream := NewStream()
	d.diversionSink.WireTo(diversionUpstream)
	return diversionStream.Writer()
}

func (d *divertToFlowStage) WireTo(stage UpstreamStage) FlowStage {
	d.upstreamStage = stage
	return d
}


// verify diversionUpstream implements UpstreamStage interface
var _ UpstreamStage = (*diversionUpstream)(nil)

type diversionUpstream struct {
	stream Stream
}

func (d *diversionUpstream) Open(ctx context.Context, mat MaterializeFunc) (Reader, *core.Future) {
	promise := core.NewPromise()
	promise.TrySuccess(NotUsed)
	return d.stream.Reader(), promise.Future()
}
