package stream

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestName(t *testing.T) {
	wg := sync.WaitGroup{}
	stream := NewStream()
	reader := stream.Reader()
	writer := stream.Writer()
		go func() {
			writer.Send(Value(1))
			wg.Add(1)
			for {
				select {
				case <- writer.Done():
					println("done")
					wg.Done()
					return
				default:

				}
			}

	}()
	element := <- reader.Read()
	reader.Complete()
	wg.Wait()
	assert.Equal(t, element.value, 1)
}
