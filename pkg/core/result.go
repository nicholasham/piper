package core

type ResultState int

const (
	IsFailure ResultState = iota
	IsSuccess
)

// Type that represents two states. Value or Error
type Result struct {
	state ResultState
	err   error
	value Any
}

type MapSuccess func(value Any) Any
type MapFailure func(err error) Any

func Success(value Any) Result {
	return Result{
		state: IsSuccess,
		err:   nil,
		value: value,
	}
}

func Failure(err error) Result {
	return Result{
		state: IsFailure,
		err:   err,
		value: nil,
	}
}

func (r Result) IsSuccess() bool {
	return r.state == IsSuccess
}

func (r Result) IsFailure() bool {
	return r.state == IsFailure
}

func (r Result) IfSuccess(f func(value Any)) Result {
	if r.IsSuccess() {
		f(r.value)
	}
	return r
}

func (r Result) Match(success MapSuccess, failure MapFailure) Any {
	if r.IsFailure() {
		return failure(r.err)
	}
	return success(r.value)
}

func (r Result) IfFailure(f func(err error)) Result {
	if r.IsFailure() {
		f(r.err)
	}
	return r
}

func (r Result) Unwrap() (Any, error) {
	if r.IsFailure() {
		return nil, r.err
	}
	return r.value, nil
}

func WrapInResult(f func() (Any, error)) Result {
	value, err := f()
	if err != nil {
		return Failure(err)
	}
	return Success(value)
}
