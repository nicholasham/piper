package stan

import (
	"github.com/nicholasham/piper/pkg/stream/sink"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestSource(t *testing.T) {

	dockerContext := newDockerContext(context.Background(), t)
	defer dockerContext.CleanUp()
	conn := dockerContext.CreateConn()
	ctx, _ := context.WithTimeout(context.Background(), 30 * time.Second)

	future := Source(conn, "test", "group1").
				To(sink.Head()).
				Run(ctx)

	conn.Publish("test", []byte("hello"))
	result := future.Await()
	value, error := result.Unwrap()

	assert.NoError(t, error)
	assert.Equal(t,  []byte("hello"), value)

}



