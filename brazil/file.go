package brazil

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
File is a wrapper around files.
*/
type File struct {
	Data *bytes.Buffer
}

/*
NewFile constructs a new file wrapper.
*/
func NewFile(path string) *File {
	buf, err := ioutil.ReadFile(
		strings.Replace(path, "~/", HomePath(), -1),
	)

	errnie.Handles(err).With(errnie.KILL)
	return &File{Data: bytes.NewBuffer(buf)}
}

/*
WriteIfNotExists is a specialized method to deal with embedded filesystems meant to
supply any missing dependencies no matter what.
*/
func WriteIfNotExists(path string, embedded embed.FS) {
	cfgFile := GetFileFromPrefix(path)
	slug := BuildPath(HomePath(), path)

	if !FileExists(slug) {
		fs := GetEmbedded(embedded, cfgFile)
		defer fs.Close()
		WriteFile(slug, ReadFile(fs))
	}
}

/*
FileExists checks if a file is present at a certain path.
*/
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

/*
GetEmbedded opens the embedded file system.
*/
func GetEmbedded(embedded embed.FS, cfgFile string) fs.File {
	fs, err := embedded.Open("cfg/" + cfgFile)
	errnie.Handles(err).With(errnie.NOOP)
	return fs
}

/*
ReadFile takes a file handle and reads the contents into a buffer.
*/
func ReadFile(fs fs.File) []byte {
	buf, err := io.ReadAll(fs)
	errnie.Handles(err).With(errnie.KILL)
	return buf
}

/*
WriteFile dumps a buffer to a file.
*/
func WriteFile(path string, buf []byte) {
	errnie.Handles(
		ioutil.WriteFile(path, buf, 0644),
	).With(errnie.NOOP)
}
