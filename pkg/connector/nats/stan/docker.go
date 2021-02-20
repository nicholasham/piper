package stan

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/net/context"
	"testing"
)

type dockerContext struct {
	t         *testing.T
	ctx       context.Context
	container testcontainers.Container
	connections [] stan.Conn
}

func newDockerContext(ctx context.Context, t *testing.T) *dockerContext {
	return &dockerContext{
		t:         t,
		ctx:       ctx,
		container: createContainer(ctx, t),
	}
}


func (d *dockerContext) CreateConn() stan.Conn {

	ip, err := d.container.Host(d.ctx)
	if err != nil {
		d.t.Error(err)
	}
	port, err := d.container.MappedPort(d.ctx, "4222")
	if err != nil {
		d.t.Error(err)
	}

	url := fmt.Sprintf("nats://%s:%d",ip, port.Int())

	conn, err := stan.Connect("test-cluster", "test", stan.NatsURL(url))
	if err != nil{
		d.t.Error(err)
	}

	d.connections = append(d.connections, conn)

	return conn
}

func (d *dockerContext) CleanUp() {
	for _, conn := range d.connections {
		conn.Close()
	}
	d.container.Terminate(d.ctx)
}

func createContainer(ctx context.Context, t *testing.T) testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image:        "nats-streaming",
		ExposedPorts: []string{"4222/tcp"},
		WaitingFor:   wait.ForLog("Streaming Server is ready"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	return container
}



