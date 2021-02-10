package stan

import (
	"context"

	"github.com/nats-io/stan.go"
	"github.com/nicholasham/piper/pkg/old-stream"
)

// verify iteratorSource implements stream.SourceStage interface
var _ old_stream.SourceStage = (*stanSourceStage)(nil)

type stanSourceStage struct {
	attributes          *old_stream.StageAttributes
	outlet              *old_stream.Outlet
	conn                stan.Conn
	subject             string
	group               string
	subscriptionOptions []stan.SubscriptionOption
}

func (s *stanSourceStage) With(options ...old_stream.StageOption) old_stream.Stage {
	attributes := s.attributes.Apply(options...)
	return &stanSourceStage{
		outlet:              old_stream.NewOutlet(attributes),
		conn:                s.conn,
		subject:             s.subject,
		group:               s.group,
		subscriptionOptions: s.subscriptionOptions,
	}
}

func (s *stanSourceStage) Name() string {
	return s.attributes.Name
}

func (s *stanSourceStage) Run(ctx context.Context) {
	go func() {
		defer s.outlet.Close()
		sub, err := s.conn.QueueSubscribe(s.subject, s.group, func(msg *stan.Msg) {
			select {
			case <-ctx.Done():
				s.outlet.SendError(ctx.Err())
				msg.Sub.Unsubscribe()
			case <-s.outlet.Done():
				msg.Sub.Unsubscribe()
			default:
			}
			s.outlet.SendValue(msg)
		}, s.subscriptionOptions...)

		if err != nil {
			s.attributes.Logger.Error(err, "failed consuming from nats")
			return
		}
		sub.Close()
	}()
}

func (s *stanSourceStage) Outlet() *old_stream.Outlet {
	return s.outlet
}

func Source(conn stan.Conn, group string, subject string, subscriptionOptions []stan.SubscriptionOption, options ...old_stream.StageOption) *old_stream.SourceGraph {
	attributes := old_stream.DefaultStageAttributes.Apply(old_stream.Name("LinearFlowStage"))
	return old_stream.FromSource(&stanSourceStage{
		outlet:              old_stream.NewOutlet(attributes),
		conn:                conn,
		subject:             subject,
		group:               group,
		subscriptionOptions: subscriptionOptions,
	})
}
