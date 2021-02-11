package fileIO

import (
	"context"
	"github.com/nicholasham/piper/pkg/stream"
	"github.com/nicholasham/piper/pkg/stream/source"
	"testing"

	"github.com/nicholasham/piper/pkg/connector/fileIO/testfs"
	"github.com/nicholasham/piper/pkg/types/list"
	"github.com/stretchr/testify/assert"
)

func TestFileSink(t *testing.T) {
	fs := testfs.New(t, "FileSinkTest")
	defer fs.CleanUp()

	expectedLines := list.OfStrings("a\n", "b\n", "c\n", "d\n", "e\n", "f\n")

	t.Run("Write lines to file", func(t *testing.T) {

		targetFile := fs.GetPath("test.txt")

		result := source.
			List(expectedLines...).
			Map(ByteString).
			RunWith(context.Background(), ToPath(targetFile, Create))

		value, err := result.Await().Unwrap()

		assert.NoError(t, err)
		assert.Equal(t, stream.NotUsed, value)
		assert.Equal(t, list.ToString(expectedLines, ""), fs.ReadFileContents(targetFile))
	})

}
