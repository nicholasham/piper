package testfs

import (
	"io/ioutil"
	"path"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

type TestFileSystem struct {
	fs                afero.Fs
	tempDirectoryPath string
	t                 *testing.T
}

func New(t *testing.T, name string) *TestFileSystem {

	fs := afero.NewOsFs()
	tempDirectoryPath, err := afero.TempDir(fs, "", name)

	if err != nil {
		t.Errorf("Failed to create test directory: %s %s", name, err)
		t.Fail()
	}

	return &TestFileSystem{
		fs:                fs,
		tempDirectoryPath: tempDirectoryPath,
		t:                 t,
	}
}

func (receiver *TestFileSystem) GetPath(fileName string) string {
	return path.Join(receiver.tempDirectoryPath, fileName)
}

func (receiver TestFileSystem) CleanUp() {
	receiver.fs.RemoveAll(receiver.tempDirectoryPath)
}

func (receiver TestFileSystem) ReadFileContents(fileName string) string {
	contents, err := ioutil.ReadFile(receiver.GetPath(strings.ReplaceAll(fileName, receiver.tempDirectoryPath, "")))
	if err != nil {
		receiver.t.Errorf("Failed to create test directory: %s %s", fileName, err)
		receiver.t.Fail()
	}
	return string(contents)
}
