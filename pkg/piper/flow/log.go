package flow

import (
	"context"
	"github.com/gammazero/workerpool"
	"github.com/nicholasham/piper/pkg/piper"
)

// verify logFlowStage implements piper.FlowStage interface
var _ piper.FlowStage = (*logFlowStage)(nil)

type logFlowStage struct {
	name       string
	attributes *piper.StageAttributes
	inlet      *piper.Inlet
	outlet     *piper.Outlet
}

func (l *logFlowStage) Name() string {
	return l.attributes.Name
}

func (l *logFlowStage) Run(ctx context.Context) {
	go func(ctx context.Context, parallelism int, logger piper.Logger, inlet *piper.Inlet, outlet *piper.Outlet) {
		wp := workerpool.New(parallelism)
		defer func() {
			outlet.Close()
		}()

		for element := range inlet.In() {

			select {
			case <-ctx.Done():
				outlet.Send(piper.Error(ctx.Err()))
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
	}(ctx, l.attributes.Parallelism, l.attributes.Logger, l.inlet, l.outlet)
}

func (l *logFlowStage) logAndSend(element piper.Element) func() {
	return func() {
		logger := l.attributes.Logger
		if !element.IsError() {
			logger.Info("[%s] value: {%v}", l.name, element.Value())
		} else {
			logger.Error(element.Error(), "[%s] Upstream failed", l.name)
		}
		l.outlet.Send(element)
	}
}

func (l *logFlowStage) Outlet() *piper.Outlet {
	return l.outlet
}

func (l *logFlowStage) Wire(stage piper.SourceStage) {
	l.inlet.WireTo(stage.Outlet())
}

func logFlow(name string, attributes ...piper.StageAttribute) piper.FlowStage {
	stageAttributes := piper.NewAttributes("HeadSink", attributes...)
	return &logFlowStage{
		name:       name,
		attributes: stageAttributes,
		inlet:      piper.NewInlet(stageAttributes),
		outlet:     piper.NewOutlet(stageAttributes),
	}
}
