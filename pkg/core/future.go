package core

type Future struct {
	resultChan chan Result
}

func (receiver *Future) Await() Result {
	result := <-receiver.resultChan
	go func() {
		receiver.resultChan <- result
	}()
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

func (receiver *Future) Then(f func(value Any) Any) *Future {
	return NewFuture(func() Result {
		result := receiver.Await()
		if result.IsSuccess() {
			return Success(f(result.value))
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
	resultChan := make(chan Result)
	go func() {
		result := f()
		resultChan <- result
	}()
	return &Future{
		resultChan: resultChan,
	}
}
