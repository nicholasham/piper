package iterator

// verify headSinkStageLogic implements SinkStageLogic interface
var _ Iterator = (*rangeIterator)(nil)

type rangeIterator struct {
	start   int
	end     int
	step    int
	current int
}

func (r *rangeIterator) ToList() []T {
	return toList(r)
}

func (r *rangeIterator) HasNext() bool {
	return r.current < (r.end * r.step)
}

func (r *rangeIterator) Next() (T, error) {
	if r.HasNext() {
		r.current = r.current + r.step
		return r.current, nil
	}
	return nil, EndOfError
}

func Range(start int, end int, step int) Iterator {
	return &rangeIterator{
		start:   start,
		end:     end,
		step:    step,
		current: 0,
	}
}
