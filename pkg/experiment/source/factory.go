package source

import (
	"github.com/nicholasham/piper/pkg/core/iterable"
	"github.com/nicholasham/piper/pkg/experiment"
)

func Single(value interface{}) *experiment.SourceGraph {
	return experiment.FromSource(experiment.SingleSource(value))
}

func Empty() *experiment.SourceGraph {
	return FromIterable(iterable.Empty())
}


func Range(start int, end int) *experiment.SourceGraph {
	return FromIterable(iterable.Range(start, end))
}

func FromIterable(iterable iterable.Iterable) *experiment.SourceGraph {
	return Single(iterable).
		MapConcat(toIterable)
}

func toIterable(value interface{}) (iterable.Iterable, error) {
	return value.(iterable.Iterable), nil
}
