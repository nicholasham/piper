package core

type Future struct {
	f func() Result
	resultChan chan Result
}

func (receiver *Future) Await() Result {
	go func() {
		receiver.resultChan <- receiver.f()
	}()
	result := <-receiver.resultChan
	return result
}

func (receiver *Future) OnSuccess(f func(value Any)) {
	go func() {
		result := receiver.Await()
		result.IfSuccess(f)
	}()
}

func (receiver *Future) OnFailure(f func(err error)) {
	go func() {
		result := receiver.Await()
		result.IfFailure(f)
	}()
}

func (receiver *Future) Then(f func(value Any) Result) *Future {
	return NewFuture(func() Result {
		result := receiver.Await()
		if result.IsSuccess() {
			return f(result.value)
		}
		return result
	})
}

func tryCompleteWith(p *Promise, f *Future) {
	go func() {
		f.Await().IfSuccess(func(value Any) {
			p.TrySuccess(value)
		}).IfFailure(func(err error) {
			p.TryFailure(err)
		})
	}()
}

func (receiver *Future) Alt(that *Future) *Future {
	p := NewPromise()
	tryCompleteWith(p, receiver)
	tryCompleteWith(p, that)
	return p.Future()
}

func NewFuture(f func() Result) *Future {
	return &Future{
		f : f,
		resultChan: make(chan Result),
	}
}
