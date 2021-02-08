package core2

// Until we get generics
type Any interface {
}


type Future interface {
	Get() Result
	OnSuccess(func(value Any) Result)
	OnFailure(func(err error) Result)
	Then(func(value Any) Result) Future
	Alt(that Future) Future
}

type Promise interface {
	Future() Future
	TrySuccess(value Any) bool
	TryFailure(err error) bool
}


// https://www.promisejs.org/implementing/

// http://www.home.hs-karlsruhe.de/~suma0002/publications/events-to-futures.pdf