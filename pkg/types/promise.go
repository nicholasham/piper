package types

import "sync"

type promiseDelivery chan Result

type Promise struct {
	sync.RWMutex
	optionalResult Optional
	waiters        []promiseDelivery
	sync.Once
}

func (p *Promise) Resolve(value interface{}) {
	p.deliver(Success(value))
}

func (p *Promise) Reject(err error) {
	p.deliver(Failure(err))
}

func (p *Promise) deliver(result Result) {
	p.Lock()
	defer p.Unlock()
	p.Once.Do(func() {
		p.optionalResult = Some(result)
		for _, w := range p.waiters {
			locW := w
			go func() {
				locW <- result
			}()
		}
	})
}

func (p *Promise) Await() Result {
	return p.optionalResult.
		Match(p.onAlreadyReceived, p.onStillWaiting).(Result)
}

func (p *Promise) onAlreadyReceived(value T) R {
	return value.(Result)
}

func (p *Promise) onStillWaiting() R {
	delivery := make(promiseDelivery)
	p.waiters = append(p.waiters, delivery)
	return <-delivery
}

func NewPromise() *Promise {
	return &Promise{
		optionalResult: None(),
		waiters:        []promiseDelivery{},
	}
}
