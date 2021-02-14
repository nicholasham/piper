package pipe

import (
	"context"
	"fmt"
	"github.com/nicholasham/piper/pkg/core"
	"github.com/nicholasham/piper/pkg/core/iterable"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"testing"
	"time"
)


func Count(ctx context.Context, outlet *Receiver) *core.Future {
	promise := core.NewPromise()

	go func(outlet *Receiver) {
		count := 0
		fmt.Println("Count started")

		for range outlet.Receive() {
			count ++
		}
		promise.TrySuccess(count)
		fmt.Println("Count finished")
	}(outlet)

	return promise.Future()
}


func Concat(ctx context.Context, upstream *Receiver) *Receiver {
	stream := NewStream("concat")
	go func(downstream *Sender, upstream *Receiver) {
		fmt.Println("Concat started")
		defer downstream.Close()
		for element := range upstream.Receive() {

			select {
			case <- ctx.Done():
				fmt.Println("Repeat timed out")
				return
			default:
			}

			select {
			case <- downstream.Done():
				fmt.Println("Concat received done signal")
				upstream.Done("concat")
			default:
			}

			element.IfOk(func(value core.Any) {
				value.(iterable.Iterable).
					ForEach(func(item interface{}) {

						//fmt.Println(fmt.Sprintf("Concat Pushing item %v", item))
						if !downstream.Send(core.Ok(item)) {
							upstream.Done("concat")
							return
						}
				})
			})

			fmt.Println("Concat receiving next")
		}

		fmt.Println("Concat finished")

	}(stream.Sender(), upstream)

	return  stream.Receiver()
}



func Repeat(ctx context.Context, value interface{}) *Receiver {
	stream := NewStream("repeat")
	go func(downstream *Sender) {
		defer downstream.Close()
		fmt.Println("Repeat started")
		for {
			select {
			case <- ctx.Done():
				fmt.Println("Repeat timed out")
				return
			default:
			}

			select {
			case <- downstream.Done():
				fmt.Println("Repeat received a signal that downstream is done receiving")
				return
			default:
			}
			if downstream.IsDone() {
				return
			}

			if !downstream.Send(core.Ok(value)) {
				return
			}

			fmt.Println("Repeating.")
		}
		fmt.Println("Repeat finished")

	}(stream.Sender())

	return stream.Receiver()
}


func Take(ctx context.Context, number int, upstreamReceiver *Receiver) *Receiver {
	stream := NewStream("take")
	go func(downstream *Sender, upstream *Receiver) {
		defer downstream.Close()
		count := 0
		fmt.Println("Take started")
		for element := range upstream.Receive() {

			select {
			case <- ctx.Done():
				fmt.Println("Repeat timed out")
				return
			default:
			}

			select {
			case <- downstream.Done():
				fmt.Println("Take received done signal")
				upstream.Done("take")
			default:
			}

			count ++
			if count <= number {
				downstream.Send(element)
				if count == number {
					upstream.Done("take")
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
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	future :=  Count(ctx,Take(ctx,10, Concat(ctx,Repeat(ctx,iterable.Range(1, 100)))))
	//future :=  Count(ctx, Take(ctx,10, Repeat(ctx, iterable.Range(1, 100))))
	result := future.Then(func(value core.Any) core.Result {
		t.Logf("Result is %v", value)
		return core.Ok(value)
	}) .Await()
	value, err := result.Unwrap()
	assert.NoError(t, err)
	assert.Equal(t, 10, value)
}