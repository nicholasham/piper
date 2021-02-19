package stan

import (
	"github.com/nats-io/stan.go"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/stream"
)

// verify publishSinkLogic implements SinkStageLogic interface
var _ stream.SinkStageLogic = (*publishSinkLogic)(nil)

type publishSinkLogic struct {
	promise  *core.Promise
	attributes          *stream.StageAttributes
	conn                stan.Conn
	subject             string
}

func (p *publishSinkLogic) OnUpstreamStart(actions stream.SinkStageActions) {
}

func (p *publishSinkLogic) OnUpstreamReceive(element stream.Element, actions stream.SinkStageActions) {
	element.WhenValue(func(value interface{}) {
		bytes, ok := value.([]byte)
		if ok {
			p.conn.Publish(p.subject, bytes )
		}
	})
}

func (p *publishSinkLogic) OnUpstreamFinish(actions stream.SinkStageActions) {
}

func publishStage(conn stan.Conn, subject string) stream.SinkStage {
	return stream.Sink(createPublishLogic(conn, subject))
}

func createPublishLogic (conn stan.Conn, subject string) stream.SinkStageLogicFactory {
	return func(attributes *stream.StageAttributes) (stream.SinkStageLogic, *core.Promise) {
		promise := core.NewPromise()
		return &publishSinkLogic{
			promise: promise,
			conn: conn,
			subject: subject,
		}, promise
	}
}

