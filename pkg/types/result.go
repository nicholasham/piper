package types

type Result struct {
	err error
	value interface{}
}

func Success(value interface{}) *Result {
	return & Result{
		err:   nil,
		value: value,
	}
}

func Failure(err error) *Result {
	return & Result{
		err:   err,
		value: nil,
	}
}

func (r *Result) Match() {

}

func (r *Result) IsSuccess() bool {
	if r.err != nil {
		return false
	}
	return true
}

func (r *Result) IsFailure() bool {
	return !r.IsSuccess()
}


func (r *Result) IfSuccess(f func(value interface{})) *Result {
	if r.IsSuccess() {
		f(r.value)
	}
	return r
}

func (r *Result) IfFailure(f func(err error)) *Result {
	if r.IsFailure() {
		f(r.err)
	}
	return r
}

