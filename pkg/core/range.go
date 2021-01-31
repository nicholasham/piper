package core

// verify headSinkStageLogic implements SinkStageLogic interface
var _ Iterator = (*rangeIterator)(nil)

type rangeIterator struct {
	start   int
	end     int
	step    int
	current int
}

func (r *rangeIterator) HasNext() bool {
	return r.current < (r.end * r.step)
}

func (r *rangeIterator) Next() interface{} {
	if r.HasNext() {
		r.current = r.current + r.step
		return r.current
	}
	return nil
}

func Range(start int, end int) Iterable {
	return SteppedRange(start, end, 1)
}

func SteppedRange(start int, end int, step int) Iterable {
	return NewIterable(func() Iterator {
		return &rangeIterator{
			start:   start,
			end:     end,
			step:    step,
			current: 0,
		}
	})
}
