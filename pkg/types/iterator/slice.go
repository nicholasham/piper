package iterator

// verify sliceIterator implements Iterator interface
var _ Iterator = (*sliceIterator)(nil)

type sliceIterator struct {
	current int
	values  []interface{}
}

func (s *sliceIterator) ToList() []interface{} {
	return toList(s)
}

func (s *sliceIterator) HasNext() bool {
	return len(s.values) >= (s.current + 1)
}

func (s *sliceIterator) Next() (interface{}, error) {
	if s.HasNext() {
		value := s.values[s.current]
		s.current = s.current + 1
		return value, nil
	}
	if len(s.values) == 0 {
		return nil, EmptyError
	}
	return nil, EndOfError
}

func Slice(values ...interface{}) Iterator {
	return &sliceIterator{
		values: values,
	}
}

func Single(value interface{}) Iterator {
	return Slice(value)
}
