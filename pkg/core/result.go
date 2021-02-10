package core

type resultState int

const (
	isFailure resultState = iota
	isSuccess
)

// Type that is used to return and propagate errors. It represents two states. Value or Error
type Result struct {
	state resultState
	err   error
	value Any
}

type MapSuccess func(value Any) Any
type MapFailure func(err error) Any

func Ok(value Any) Result {
	return Result{
		state: isSuccess,
		err:   nil,
		value: value,
	}
}

func Err(err error) Result {
	return Result{
		state: isFailure,
		err:   err,
		value: nil,
	}
}

func (r Result) IsOk() bool {
	return r.state == isSuccess
}

func (r Result) IsErr() bool {
	return r.state == isFailure
}

func (r Result) Map(f MapSuccess) Result {
	if r.IsOk() {
		return Ok(f(r.value))
	}
	return r
}

func (r Result) Then(f func(value Any) Result) Result {
	if r.IsOk() {
		return f(r.value)
	}
	return r
}

func (r Result) OrElse(f func(err error) Result) Result {
	if r.IsErr() {
		return f(r.err)
	}
	return r
}

func (r Result) IfOk(f func(value Any)) Result {
	if r.IsOk() {
		f(r.value)
	}
	return r
}

func (r Result) Match(success MapSuccess, failure MapFailure) Any {
	if r.IsErr() {
		return failure(r.err)
	}
	return success(r.value)
}

func (r Result) IfErr(f func(err error)) Result {
	if r.IsErr() {
		f(r.err)
	}
	return r
}

func (r Result) Unwrap() (Any, error) {
	if r.IsErr() {
		return nil, r.err
	}
	return r.value, nil
}
