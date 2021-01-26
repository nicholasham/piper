package stan

import (
	"context"

	"github.com/nats-io/stan.go"
	"github.com/nicholasham/piper/pkg/stream"
)

// verify iteratorSource implements stream.SourceStage interface
var _ stream.SourceStage = (*stanSourceStage)(nil)

type stanSourceStage struct {
	name                string
	logger              stream.Logger
	outlet              *stream.Outlet
	conn                stan.Conn
	subject             string
	group               string
	subscriptionOptions []stan.SubscriptionOption
}

func (s *stanSourceStage) Name() string {
	return s.name
}

func (s *stanSourceStage) Run(ctx context.Context) {
	go func() {
		defer s.outlet.Close()
		sub, err := s.conn.QueueSubscribe(s.subject, s.group, func(msg *stan.Msg) {
			select {
			case <-ctx.Done():
				s.outlet.Send(stream.Error(ctx.Err()))
				msg.Sub.Unsubscribe()
			case <-s.outlet.Done():
				msg.Sub.Unsubscribe()
			default:
			}
			s.outlet.Send(stream.Value(msg))
		}, s.subscriptionOptions...)

		if err != nil {
			s.logger.Error(err, "failed consuming from nats")
			return
		}
		sub.Close()
	}()
}

func (s *stanSourceStage) Outlet() *stream.Outlet {
	return s.outlet
}

func Source(conn stan.Conn, group string, subject string, subscriptionOptions []stan.SubscriptionOption, options ... stream.StageOption) *stream.SourceGraph {
	state := stream.NewStageState("StanSource", options...)
	return stream.SourceFrom(&stanSourceStage{
		name:                state.Name,
		logger:              state.Logger,
		outlet:              stream.NewOutlet(state),
		conn:                conn,
		subject:             subject,
		group:               group,
		subscriptionOptions: subscriptionOptions,
	})
}
