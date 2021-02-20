package iterable

// verify rangeIterator implements Iterator interface
var _ Iterator = (*rangeIterator)(nil)

type rangeIterator struct {
	end     int
	step    int
	current int
}

func (r *rangeIterator) HasNext() bool {
	return r.current <= (r.end * r.step)
}

func (r *rangeIterator) Next() interface{} {
	if r.HasNext() {
		item := r.current
		r.current = r.current + r.step
		return item
	}
	panic("next was called when has no next value, always check there is a next value by calling HasNext.")
}

func Range(start int, end int) Iterable {
	return SteppedRange(start, end, 1)
}

func SteppedRange(start int, end int, step int) Iterable {
	return NewIterable(func() Iterator {
		return &rangeIterator{
			end:     end,
			step:    step,
			current: start,
		}
	})
}
