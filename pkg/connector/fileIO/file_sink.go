package fileIO

import (
	"fmt"
	"github.com/nicholasham/piper/pkg/core"
	"os"

	"github.com/nicholasham/piper/pkg/stream"
)

// verify fileSinkStageLogic implements SinkStageLogic interface
var _ stream.SinkStageLogic = (*fileSinkStageLogic)(nil)

var ByteArrayError = fmt.Errorf("expected element value to be a byte array")

type fileSinkStageLogic struct {
	promise  *core.Promise
	filePath string
	file     *os.File
	factory  FileFactory
}

func (f *fileSinkStageLogic) OnUpstreamStart(actions stream.SinkStageActions) {
	file, err := f.factory(f.filePath)
	if err != nil {
		actions.FailStage(err)
	}
	f.file = file
}

func (f *fileSinkStageLogic) OnUpstreamReceive(element stream.Element, actions stream.SinkStageActions) {
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

func (f *fileSinkStageLogic) OnUpstreamFinish(actions stream.SinkStageActions) {
	f.file.Close()
	f.promise.TrySuccess(stream.NotUsed)
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

func createFileSinkLogic(filePath string, factory FileFactory) stream.SinkStageLogicFactory {
	return func(attributes *stream.StageAttributes) (stream.SinkStageLogic, *core.Promise) {
		promise := core.NewPromise()
		return &fileSinkStageLogic{
			promise:  promise,
			filePath: filePath,
			factory:  factory,
		}, promise
	}
}

func ToPath(filePath string, factory FileFactory) *stream.SinkGraph {
	return stream.FromSink(stream.Sink(createFileSinkLogic(filePath, factory)))
}
