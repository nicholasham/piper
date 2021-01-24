package io

import (
	"context"
	"fmt"
	"os"

	"github.com/nicholasham/piper/pkg/stream"
	"github.com/nicholasham/piper/pkg/stream/sink"
)

// verify fileCollector implements CollectorLogic interface
var _ sink.CollectorLogic = (*fileCollector)(nil)

var ByteArrayError = fmt.Errorf("expected element value to be a byte array")

type fileCollector struct {
	filePath string
	file     *os.File
	factory  FileFactory
}

type FileFactory func(filePath string) (*os.File, error)

func Create(filePath string) (*os.File, error) {
	return os.Create(filePath)
}

func Append(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, os.O_APPEND, os.ModeAppend)
}

func Custom(flag int, perm os.FileMode) FileFactory {
	return func(filePath string) (*os.File, error) {
		return os.OpenFile(filePath, flag, perm)
	}
}

func (f *fileCollector) Start(ctx context.Context, actions sink.CollectActions) {
	file, err := f.factory(f.filePath)
	if err != nil {
		actions.FailStage(err)
	}
	f.file = file
}

func (f *fileCollector) Collect(ctx context.Context, element stream.Element, actions sink.CollectActions) {
	element.WhenValue(func(value interface{}) {
		bytes, ok := value.([]byte)
		if !ok {
			actions.FailStage(ByteArrayError)
		}
		_, err := f.file.Write(bytes)
		if err != nil {
			actions.FailStage(err)
		}
	})
	element.WhenError(actions.FailStage)
}

func (f *fileCollector) End(ctx context.Context, actions sink.CollectActions) {
	f.file.Close()
}

func Sink(filePath string, factory FileFactory) *stream.SinkGraph {
	return sink.Collector("FileSink",
		&fileCollector{
			filePath: filePath,
			factory:  factory,
		})
}
