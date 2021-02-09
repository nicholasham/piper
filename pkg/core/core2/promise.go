package core2

// Until we get generics
type Any interface {
}

type Future struct {
	resultChan chan Result
}

func (receiver *Future) Get() Result {
	result := <- receiver.resultChan
	go func() {
		receiver.resultChan <- result
	}()
	return result
}

func (receiver *Future) OnSuccess(f func(value Any)) {
	go func() {
		result := receiver.Get()
		result.IfSuccess(f)
	}()
}

func (receiver *Future) OnFailure(f func(err error)) {
	go func() {
		result := receiver.Get()
		result.IfFailure(f)
	}()
}

func (receiver *Future) Then(f func(value Any) Result) *Future {
	return NewFuture(func() Result {
		result := receiver.Get()
		if result.IsSuccess() {
			return f(result.value)
		}
		return result
	})
}

func tryCompleteWith(p *Promise, f *Future) {
	go func() {
		f.Get().IfSuccess(func(value Any) {
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

func NewFuture(f func() Result) *Future  {
	resultChan := make(chan Result)
	go func() {
		result := f()
		resultChan <- result
	}()
	return &Future{
		resultChan: resultChan,
	}
}

type Promise struct {
	completed bool
	resultChan chan Result
}

func (p *Promise) Future() *Future {
	return NewFuture(func() Result {
		result := <- p.resultChan
		return result
	})
}

func (p *Promise) TrySuccess(value Any) bool {
	if !p.completed {
		go func() {
			p.resultChan <- Success(value)
		}()
		return true
	}
	return false
}

func (p *Promise) TryFailure(err error) bool {
	if !p.completed {
		go func() {
			p.resultChan <- Failure(err)
		}()
		return true
	}
	return false
}

func NewPromise() *Promise {
	return &Promise{
		completed: false,
		resultChan: make(chan Result),
	}
}

// https://www.promisejs.org/implementing/

// http://www.home.hs-karlsruhe.de/~suma0002/publications/events-to-futures.pdf