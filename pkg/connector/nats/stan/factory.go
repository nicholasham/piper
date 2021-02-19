package stan

import (
	"github.com/nats-io/stan.go"
	"github.com/nicholasham/piper/pkg/stream"
)

func PublishSink(conn stan.Conn, subject string) * stream.SinkGraph {
	return stream.FromSink(publishStage(conn, subject))
}
