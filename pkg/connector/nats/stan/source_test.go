package stan

import (
	"context"
	"github.com/nats-io/stan.go"
	"github.com/nicholasham/piper/pkg/stream/sink"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSource(t *testing.T) {

	dockerContext := newDockerContext(context.Background(), t)
	defer dockerContext.CleanUp()
	conn := dockerContext.CreateConn()
	ctx, _ := context.WithTimeout(context.Background(), 50 * time.Second)

	future := Source(conn, "test", "group1", stan.StartAt(0)).
				Map(func(value interface{}) (interface{}, error) {
					return value.(*stan.Msg).Data, nil
				}).
				To(sink.Head()).
				Run(ctx)

	conn.Publish("test", []byte("hello"))

	result := future.Await()
	value, error := result.Unwrap()

	assert.NoError(t, error)
	assert.Equal(t,  []byte("hello"), value)

}



