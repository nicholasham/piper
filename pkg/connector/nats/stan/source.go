package stan

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/stream"

	"github.com/nats-io/stan.go"
)

// verify iteratorSource implements stream.SourceStage interface
var _ stream.SourceStage = (*stanSourceStage)(nil)

type stanSourceStage struct {
	attributes          *stream.StageAttributes
	conn                stan.Conn
	subject             string
	group               string
	subscriptionOptions []stan.SubscriptionOption
}

func (s *stanSourceStage) Open(ctx context.Context, mat stream.MaterializeFunc) (stream.Reader, *core.Future) {
	logger := s.attributes.Logger
	outputPromise := core.NewPromise()
	outputStream := stream.NewStream(s.attributes.Name)
	go func() {
		writer := outputStream.Writer()
		defer writer.Close()

		sub, err := s.conn.QueueSubscribe(s.subject, s.group, func(msg *stan.Msg) {
			select {
			case <-ctx.Done():
				writer.SendError(ctx.Err())
				msg.Sub.Unsubscribe()
			case <-writer.Done():
				msg.Sub.Unsubscribe()
			default:
			}
			writer.SendValue(msg)
		}, s.subscriptionOptions...)

		if err != nil {
			logger.Error(err, "failed consuming from nats")
			return
		}
		sub.Close()
	}()
	return outputStream.Reader(), outputPromise.Future()
}

func (s *stanSourceStage) With(options ...stream.StageOption) stream.Stage {
	return &stanSourceStage{
		attributes:          s.attributes.Apply(options...),
		conn:                s.conn,
		subject:             s.subject,
		group:               s.group,
		subscriptionOptions: s.subscriptionOptions,
	}
}

func Source(conn stan.Conn, group string, subject string, subscriptionOptions []stan.SubscriptionOption) *stream.SourceGraph {
	return stream.FromSource(&stanSourceStage{
		attributes:          stream.DefaultStageAttributes.Apply(stream.Name("StanSource")),
		conn:                conn,
		subject:             subject,
		group:               group,
		subscriptionOptions: subscriptionOptions,
	})
}
