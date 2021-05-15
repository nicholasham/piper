package pipe

import (
	"context"
	"fmt"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/core/iterable"
	"github.com/nicholasham/piper/pkg/stream"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"sync"
	"testing"
	"time"
)


func Count2(ctx context.Context, done chan struct{}, upstream chan stream.Element) chan stream.Element {

}

func Count(ctx context.Context, wg *sync.WaitGroup, outlet *Receiver) *core.Future {
	promise := core.NewPromise()
	wg.Add(1)
	go func(outlet *Receiver) {
		defer wg.Done()
		count := 0
		fmt.Println("Count started")
		for range outlet.Receive() {
			count++
		}
		promise.TrySuccess(count)
		fmt.Println("Count finished")
	}(outlet)

	return promise.Future()
}

func cleanUp(funcs ... func()) {
	for _, f := range funcs {
		f()
	}
}

func Concat(ctx context.Context, wg *sync.WaitGroup, upstream *Receiver) *Receiver {
	stream := NewStream("concat")
	wg.Add(1)
	go func(downstream *Sender, upstream *Receiver) {
		fmt.Println("Concat started")
		defer cleanUp(downstream.Close, wg.Done)
		for element := range upstream.Receive() {

			select {
			case <-ctx.Done():
				fmt.Println("Repeat timed out")
				return
			default:
			}

			select {
			case <-downstream.Done():
				fmt.Println("Concat received done signal")
				upstream.Done()
				return
			default:
			}

			element.IfOk(func(value core.Any) {
				value.(iterable.Iterable).
					TakeWhile(func(value core.Any) bool {
						if downstream.IsDone() {
							fmt.Println("concat downstream done")
							upstream.Done()
							return false
						}
						return true
					}).
					ForEach(func(item interface{}) {
						fmt.Println("concat sending ...")
						downstream.TrySend(core.Ok(item))
					})
			})

			if downstream.IsDone() {
				break
			}

			fmt.Println("Concat receiving next")
		}

		fmt.Println("Concat finished")

	}(stream.Sender(), upstream)

	return stream.Receiver()
}

func Repeat(ctx context.Context, wg *sync.WaitGroup, value interface{}) *Receiver {
	stream := NewStream("repeat")
	wg.Add(1)
	go func(downstream *Sender) {
		defer cleanUp(downstream.Close, wg.Done)
		fmt.Println("Repeat started")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Repeat timed out")
				return
			default:
			}

			select {
			case <-downstream.Done():
				fmt.Println("Repeat received a signal that downstream is done receiving")
				return
			default:
			}

			if downstream.IsDone() {
				return
			}

			if !downstream.TrySend(core.Ok(value)) {
				return
			}

			fmt.Println("Repeating.")
		}
		fmt.Println("Repeat finished")

	}(stream.Sender())

	return stream.Receiver()
}

func Take(ctx context.Context, wg *sync.WaitGroup, number int, upstreamReceiver *Receiver) *Receiver {
	stream := NewStream("take")
	wg.Add(1)
	go func(downstream *Sender, upstream *Receiver) {
		defer cleanUp(downstream.Close, wg.Done)
		count := 0
		fmt.Println("Take started")
		for element := range upstream.Receive() {

			select {
			case <-ctx.Done():
				fmt.Println("Repeat timed out")
				return
			default:
			}

			select {
			case <-downstream.Done():
				fmt.Println("Take received done signal")
				upstream.Done()
				return
			default:
			}

			count++
			if count <= number {
				downstream.TrySend(element)
				if count == number {
					upstream.Done()
					break
				}
			}

			fmt.Println("Take receiving next")

		}
		fmt.Println("Take finished")

	}(stream.Sender(), upstreamReceiver)
	return stream.Receiver()
}

func TestUpstreamIsClosed(t *testing.T) {
	goleak.VerifyNone(t)
	wg := &sync.WaitGroup{}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	future := Count(ctx,wg, Take(ctx, wg,10, Concat(ctx,wg, Repeat(ctx,wg, iterable.Range(1, 100)))))
	//future :=  Count(ctx, Take(ctx,10, Repeat(ctx, iterable.Range(1, 100))))
	result := future.Then(func(value core.Any) core.Result {
		t.Logf("Result is %v", value)
		return core.Ok(value)
	}).Await()
	wg.Wait()
	value, err := result.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, 10, value)
}
