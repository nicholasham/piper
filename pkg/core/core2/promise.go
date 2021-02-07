package core2

import "sync"

// Until we get generics
type Any interface {
}

type Result interface {

}

type Future interface {
	Await()
}

type Promise struct {
	wg  sync.WaitGroup
}

func NewPromise(resolve func(value Any), reject func(err error)) * Promise {
	return &Promise{}
}

func (p *Promise) Map(func(value Any) Any)  *Promise {
	return &Promise{

	}
}


// https://www.promisejs.org/implementing/

// http://www.home.hs-karlsruhe.de/~suma0002/publications/events-to-futures.pdf