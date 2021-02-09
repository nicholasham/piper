package core

// Until we get generics
type Any interface {
}

type Promise struct {
	completed  bool
	resultChan chan Result
}

func (p *Promise) Future() *Future {
	return NewFuture(func() Result {
		result := <-p.resultChan
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
		completed:  false,
		resultChan: make(chan Result),
	}
}

// http://www.home.hs-karlsruhe.de/~suma0002/publications/events-to-futures.pdf
