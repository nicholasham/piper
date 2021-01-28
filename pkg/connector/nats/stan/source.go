package stan

import (
	"context"

	"github.com/nats-io/stan.go"
	"github.com/nicholasham/piper/pkg/streamold"
)

// verify iteratorSource implements stream.SourceStage interface
var _ streamold.SourceStage = (*stanSourceStage)(nil)

type stanSourceStage struct {
	name                string
	logger              streamold.Logger
	outlet              *streamold.Outlet
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
				s.outlet.Send(streamold.Error(ctx.Err()))
				msg.Sub.Unsubscribe()
			case <-s.outlet.Done():
				msg.Sub.Unsubscribe()
			default:
			}
			s.outlet.Send(streamold.Value(msg))
		}, s.subscriptionOptions...)

		if err != nil {
			s.logger.Error(err, "failed consuming from nats")
			return
		}
		sub.Close()
	}()
}

func (s *stanSourceStage) Outlet() *streamold.Outlet {
	return s.outlet
}

func Source(conn stan.Conn, group string, subject string, subscriptionOptions []stan.SubscriptionOption, options ... streamold.StageOption) *streamold.SourceGraph {
	stageOptions := streamold.DefaultStageOptions.Apply(streamold.Name("StanSource")).Apply(options...)
	return streamold.SourceFrom(&stanSourceStage{
		name:                stageOptions.Name,
		logger:              stageOptions.Logger,
		outlet:              streamold.NewOutlet(stageOptions),
		conn:                conn,
		subject:             subject,
		group:               group,
		subscriptionOptions: subscriptionOptions,
	})
}
