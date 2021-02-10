package source

import (
	"github.com/nicholasham/piper/pkg/core/iterable"
	"github.com/nicholasham/piper/pkg/stream"
)

func Single(value interface{}) *stream.SourceGraph {
	return stream.FromSource(stream.SingleStage(value))
}

func Empty() *stream.SourceGraph {
	return FromIterable(iterable.Empty())
}


func Range(start int, end int) *stream.SourceGraph {
	return FromIterable(iterable.Range(start, end))
}

func FromIterable(iterable iterable.Iterable) *stream.SourceGraph {
	return Single(iterable).
		MapConcat(toIterable)
}

func toIterable(value interface{}) (iterable.Iterable, error) {
	return value.(iterable.Iterable), nil
}
