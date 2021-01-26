package flow

import (
	"context"

	"github.com/gammazero/workerpool"
	"github.com/nicholasham/piper/pkg/stream"
)

// verify logFlowStage implements stream.FlowStage interface
var _ stream.FlowStage = (*logFlowStage)(nil)

type logFlowStage struct {
	name string
	logger stream.Logger
	parallelism int
	inlet      *stream.Inlet
	outlet     *stream.Outlet
}

func (l *logFlowStage) Name() string {
	return l.name
}

func (l *logFlowStage) Run(ctx context.Context) {
	go func(ctx context.Context, parallelism int, logger stream.Logger, inlet *stream.Inlet, outlet *stream.Outlet) {
		wp := workerpool.New(parallelism)
		defer func() {
			outlet.Close()
		}()

		for element := range inlet.In() {

			select {
			case <-ctx.Done():
				outlet.Send(stream.Error(ctx.Err()))
				inlet.Complete()
			case <-outlet.Done():
				inlet.Complete()
			default:
			}

			if !inlet.CompletionSignaled() {
				wp.Submit(l.logAndSend(element))
			}
		}

		wp.StopWait()
	}(ctx, l.parallelism, l.logger, l.inlet, l.outlet)
}

func (l *logFlowStage) logAndSend(element stream.Element) func() {
	return func() {
		logger := l.logger
		if !element.IsError() {
			logger.Info("[%s] value: {%v}", l.name, element.Value())
		} else {
			logger.Error(element.Error(), "[%s] Upstream failed", l.name)
		}
		l.outlet.Send(element)
	}
}

func (l *logFlowStage) Outlet() *stream.Outlet {
	return l.outlet
}

func (l *logFlowStage) Wire(stage stream.SourceStage) {
	l.inlet.WireTo(stage.Outlet())
}

func logFlow(name string, options ...stream.StageOption) stream.FlowStage {
	stageOptions := stream.DefaultStageOptions.
		Apply(stream.Name(name)).
		Apply(options...)

	return &logFlowStage{
		name: stageOptions.Name,
		logger: stageOptions.Logger,
		parallelism: stageOptions.Parallelism,
		inlet:      stream.NewInlet(stageOptions),
		outlet:     stream.NewOutlet(stageOptions),
	}
}
