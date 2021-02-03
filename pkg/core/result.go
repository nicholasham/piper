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
	value A
}

// Represents the value type
type A interface {
}

// Represents result type
type R interface {
}

type MapSuccess func(value A) R
type MapFailure func(err error) R

func Success(value A) Result {
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

func (r Result) IfSuccess(f func(value A)) Result {
	if r.IsSuccess() {
		f(r.value)
	}
	return r
}

func (r Result) Match(success MapSuccess, failure MapFailure) R {
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

func (r Result) Unwrap() (A, error) {
	if r.IsFailure() {
		return nil, r.err
	}
	return r.value, nil
}

func WrapInResult(f func() (A, error)) Result {
	value, err := f()
	if err != nil {
		return Failure(err)
	}
	return Success(value)
}
