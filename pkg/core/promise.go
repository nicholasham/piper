package core

import "sync"

// Until we get generics
type Any interface {
}

type Promise struct {
	completed  bool
	resultChan chan Result
	mu sync.RWMutex
}

func (p *Promise) Future() *Future {
	return NewFuture(func() Result {
		result := <-p.resultChan
		return result
	})
}

func (p *Promise) IsCompleted() bool {
	return p.completed
}

func (p *Promise) TrySuccess(value Any) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.completed {
		go func() {
			p.resultChan <- Success(value)
		}()
		p.completed = true
		return true
	}
	return false
}

func (p *Promise) TryFailure(err error) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.completed {
		go func() {
			p.resultChan <- Failure(err)
		}()
		p.completed = true
		return true
	}
	return false
}

func NewPromise() *Promise {
	return &Promise{
		completed:  false,
		resultChan: make(chan Result, 1),
	}
}

// http://www.home.hs-karlsruhe.de/~suma0002/publications/events-to-futures.pdf
