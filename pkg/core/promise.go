package core

import "sync"

// http://www.home.hs-karlsruhe.de/~suma0002/publications/events-to-futures.pdf

type Promise struct {
	completed  bool
	resultChan chan Result
	mu         sync.RWMutex
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
			p.resultChan <- Ok(value)
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
			p.resultChan <- Err(err)
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
