package flow

import "sync"

type promiseDelivery chan *result
type result struct {
	value interface{}
	err   error
}
type Promise struct {
	sync.RWMutex
	result    *result
	waiters   []promiseDelivery
	completed bool
	sync.Once
}

func (p *Promise) Resolve(value interface{}) {
	p.deliver(&result{
		value: value,
	})
}

func (p *Promise) Reject(err error) {
	p.deliver(&result{
		err: err,
	})
}

func (p *Promise) deliver(result *result) {
	p.Lock()
	defer p.Unlock()
	p.Once.Do(func() {
		if !p.completed {
			p.result = result
			p.completed = true
		}
		for _, w := range p.waiters {
			locW := w
			go func() {
				locW <- result
			}()
		}
	})
}

func (p *Promise) Await() (interface{}, error) {
	if p.result != nil {
		return p.result.value, p.result.err
	}
	delivery := make(promiseDelivery)
	p.waiters = append(p.waiters, delivery)
	result := <-delivery
	return result.value, result.err
}

func NewPromise() *Promise {
	return &Promise{
		result:  nil,
		waiters: []promiseDelivery{},
	}
}
