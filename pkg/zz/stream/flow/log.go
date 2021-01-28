package flow

import (
	"context"

	"github.com/gammazero/workerpool"
	"github.com/nicholasham/piper/pkg/streamold"
)

// verify logFlowStage implements stream.FlowStage interface
var _ streamold.FlowStage = (*logFlowStage)(nil)

type logFlowStage struct {
	name        string
	logger      streamold.Logger
	parallelism int
	inlet       *streamold.Inlet
	outlet      *streamold.Outlet
}

func (l *logFlowStage) Name() string {
	return l.name
}

func (l *logFlowStage) Run(ctx context.Context) {
	go func(ctx context.Context, parallelism int, logger streamold.Logger, inlet *streamold.Inlet, outlet *streamold.Outlet) {
		wp := workerpool.New(parallelism)
		defer func() {
			outlet.Close()
		}()

		for element := range inlet.In() {

			select {
			case <-ctx.Done():
				outlet.Send(streamold.Error(ctx.Err()))
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

func (l *logFlowStage) logAndSend(element streamold.Element) func() {
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

func (l *logFlowStage) Outlet() *streamold.Outlet {
	return l.outlet
}

func (l *logFlowStage) Wire(stage streamold.SourceStage) {
	l.inlet.WireTo(stage.Outlet())
}

func logFlow(name string, options ...streamold.StageOption) streamold.FlowStage {
	stageOptions := streamold.DefaultStageOptions.
		Apply(streamold.Name(name)).
		Apply(options...)

	return &logFlowStage{
		name:        stageOptions.Name,
		logger:      stageOptions.Logger,
		parallelism: stageOptions.Parallelism,
		inlet:       streamold.NewInlet(stageOptions),
		outlet:      streamold.NewOutlet(stageOptions),
	}
}
