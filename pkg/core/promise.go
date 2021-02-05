package core

import "sync"

type promiseDelivery chan Result

type Promise struct {
	mutex sync.RWMutex
	optionalResult Optional
	waiters        []promiseDelivery
	once sync.Once
	f func(value interface{}) interface{}
}

func (p *Promise) TrySuccess(value interface{}) {
	if p.f != nil {
		p.deliver(Success(value))
	}else{
		p.deliver(Success(value))
	}
}

func (p *Promise) TryFailure(err error) {
	p.deliver(Failure(err))
}

func (p *Promise) deliver(result Result) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.once.Do(func() {
		p.optionalResult = Some(result)
		for _, w := range p.waiters {
			locW := w
			go func() {
				locW <- result
			}()
		}
	})
}

func (p *Promise) FlatMap(f func(value interface{}) interface{}) *Promise  {
	return p
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
