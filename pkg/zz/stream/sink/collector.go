package sink

import (
	"context"

	"github.com/nicholasham/piper/pkg/streamold"
)

type CollectorLogic interface {
	Start(ctx context.Context, actions CollectActions)
	Collect(ctx context.Context, element streamold.Element, actions CollectActions)
	End(ctx context.Context, actions CollectActions)
}

type CollectActions interface {
	FailStage(cause error)
	CompleteStage(value interface{})
}

type StartActions interface {
	FailStage(cause error)
	CompleteStage(value interface{})
}

type EndActions interface {
	FailStage(cause error)
	CompleteStage(value interface{})
}

// verify collectActions implements stream.CollectActions interface
var _ CollectActions = (*collectActions)(nil)

type collectActions struct {
	failStage     func(cause error)
	completeStage func(value interface{})
}

func (c *collectActions) FailStage(cause error) {
	c.failStage(cause)
}

func (c *collectActions) CompleteStage(value interface{}) {
	c.completeStage(value)
}

// verify collectorSinkStage implements stream.SinkStage interface
var _ streamold.SinkStage = (*collectorSinkStage)(nil)

type collectorSinkStage struct {
	name      string
	logger    streamold.Logger
	collector CollectorLogic
	inlet     *streamold.Inlet
	promise   *streamold.Promise
}

func (c *collectorSinkStage) Name() string {
	return c.name
}

func (c *collectorSinkStage) newActions() CollectActions {
	return &collectActions{
		failStage: func(cause error) {
			c.logger.Error(cause, "failed stage because")
			c.inlet.Complete()
			c.promise.Reject(cause)
		},
		completeStage: func(value interface{}) {
			c.inlet.Complete()
			c.promise.Resolve(value)
		},
	}
}

func (c *collectorSinkStage) Run(ctx context.Context) {
	go func() {
		actions := c.newActions()
		c.collector.Start(ctx, actions)
		for element := range c.inlet.In() {
			select {
			case <-ctx.Done():
				c.promise.Reject(ctx.Err())
				c.inlet.Complete()
			default:

			}
			if !c.inlet.CompletionSignaled() {
				c.collector.Collect(ctx, element, c.newActions())
			}
		}
		c.collector.End(ctx, c.newActions())
		c.promise.Resolve(streamold.NotUsed)
	}()
}

func (c *collectorSinkStage) Wire(stage streamold.SourceStage) {
	c.inlet.WireTo(stage.Outlet())
}

func (c *collectorSinkStage) Inlet() *streamold.Inlet {
	return c.inlet
}

func (c *collectorSinkStage) Result() streamold.Future {
	return c.promise
}

func CollectorSink(name string, collector CollectorLogic, options []streamold.StageOption) streamold.SinkStage {
	stageOptions := streamold.DefaultStageOptions.
		Apply(streamold.Name(name)).
		Apply(options...)
	return &collectorSinkStage{
		name:      stageOptions.Name,
		logger:    stageOptions.Logger,
		collector: collector,
		inlet:     streamold.NewInlet(stageOptions),
		promise:   streamold.NewPromise(),
	}
}
