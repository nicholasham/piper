package core

// verify sliceIterator implements Iterator interface
var _ Iterator = (*sliceIterator)(nil)

type sliceIterator struct {
	current int
	values  []interface{}
}

func (s *sliceIterator) HasNext() bool {
	return len(s.values) >= (s.current + 1)
}

func (s *sliceIterator) Next() interface{} {
	if s.HasNext() {
		value := s.values[s.current]
		s.current = s.current + 1
		return value
	}
	return nil
}

func Slice(values ...interface{}) Iterable {
	return NewIterable(func() Iterator {
		return &sliceIterator{
			values: values,
		}
	})

}
