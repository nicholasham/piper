package core2

// Until we get generics
type Any interface {
}


type Future interface {
	Get() Result
	OnSuccess(func(value Any))
	OnFailure(func(err error))
	Then(func(value Any) Result) Future
	Alt(that Future) Future
}

// verify future implements Future interface
var _ Future = (*future)(nil)

type future struct {
	resultChan chan Result
}

func (receiver *future) Get() Result {
	result := <- receiver.resultChan
	go func() {
		receiver.resultChan <- result
	}()
	return result
}

func (receiver *future) OnSuccess(f func(value Any)) {
	go func() {
		result := receiver.Get()
		result.IfSuccess(f)
	}()
}

func (receiver *future) OnFailure(f func(err error)) {
	go func() {
		result := receiver.Get()
		result.IfFailure(f)
	}()
}

func (receiver *future) Then(f func(value Any) Result) Future {
	return NewFuture(func() Result {
		result := receiver.Get()
		if result.IsSuccess() {
			return f(result.value)
		}
		return result
	})
}

func tryCompleteWith(p Promise, f Future) {
	go func() {
		f.Get().IfSuccess(func(value Any) {
			p.TrySuccess(value)
		}).IfFailure(func(err error) {
			p.TryFailure(err)
		})
	}()
}

func (receiver *future) Alt(that Future) Future {
	p := NewPromise()
	tryCompleteWith(p, receiver)
	tryCompleteWith(p, that)
	return p.Future()
}

func NewFuture(f func() Result) Future  {
	resultChan := make(chan Result)
	go func() {
		result := f()
		resultChan <- result
	}()
	return & future{
		resultChan: resultChan,
	}
}

type Promise interface {
	Future() Future
	TrySuccess(value Any) bool
	TryFailure(err error) bool
}

// verify promise implements Promise interface
var _ Promise = (*promise)(nil)


type promise struct {
	completed bool
	resultChan chan Result
}

func (p *promise) Future() Future {
	return NewFuture(func() Result {
		result := <- p.resultChan
		return result
	})
}

func (p *promise) TrySuccess(value Any) bool {
	if !p.completed {
		go func() {
			p.resultChan <- Success(value)
		}()
		return true
	}
	return false
}

func (p *promise) TryFailure(err error) bool {
	if !p.completed {
		go func() {
			p.resultChan <- Failure(err)
		}()
		return true
	}
	return false
}

func NewPromise() Promise {
	return &promise {
		completed: false,
		resultChan: make(chan Result),
	}
}

// https://www.promisejs.org/implementing/

// http://www.home.hs-karlsruhe.de/~suma0002/publications/events-to-futures.pdf