package stan

import (
	"context"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/stream"
	"sync"

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

func (s *stanSourceStage) Named(name string) stream.Stage {
	return s.With(stream.Name(name))
}

func (s *stanSourceStage) Open(ctx context.Context, wg *sync.WaitGroup, mat stream.MaterializeFunc) (*stream.Receiver, *core.Future) {
	logger := s.attributes.Logger
	outputPromise := core.NewPromise()
	outputStream := stream.NewStream(s.attributes.Name)
	go func() {
		sender := outputStream.Sender()
		defer sender.Close()
		wg.Add(1)

		sub, err := s.conn.QueueSubscribe(s.subject, s.group, func(msg *stan.Msg) {
			sender.TrySend(stream.Value(msg))
		}, s.subscriptionOptions...)

		if err != nil {
			logger.Error(err, "failed subscribing")
			return
		}

		go func() {
			for {
				select {
				case <-ctx.Done():
					sender.TrySend(stream.Error(ctx.Err()))
					wg.Done()
					sub.Unsubscribe()
					return
				case <-sender.Done():
					wg.Done()
					sub.Unsubscribe()
					return
				default:
				}
			}
		}()

		wg.Wait()
		sub.Close()
	}()
	return outputStream.Receiver(), outputPromise.Future()
}

func (s *stanSourceStage) With(options ...stream.StageOption) stream.Stage {
	return &stanSourceStage{
		attributes:          s.attributes.With(options...),
		conn:                s.conn,
		subject:             s.subject,
		group:               s.group,
		subscriptionOptions: s.subscriptionOptions,
	}
}

func Source(conn stan.Conn, subject string, group string, subscriptionOptions ...stan.SubscriptionOption) *stream.SourceGraph {
	return stream.FromSource(&stanSourceStage{
		attributes:          stream.DefaultStageAttributes.With(stream.Name("StanSource")),
		conn:                conn,
		subject:             subject,
		group:               group,
		subscriptionOptions: subscriptionOptions,
	})
}
