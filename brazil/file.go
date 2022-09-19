package brazil

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
File is a wrapper around files.
*/
type File struct {
	Data bytes.Buffer
}

/*
NewFile constructs a new file wrapper.
*/
func NewFile(path string) File {
	errnie.Traces()
	buf, err := os.ReadFile(
		strings.ReplaceAll(path, "~", HomePath()),
	)

	errnie.Handles(err)
	return File{Data: *bytes.NewBuffer(buf)}
}

/*
WriteIfNotExists is a specialized method to deal with embedded filesystems meant to
supply any missing dependencies no matter what.
*/
func WriteIfNotExists(path string, embedded embed.FS, ex bool) {
	errnie.Traces()

	if !FileExists(path) {
		fs := GetEmbedded(embedded, path)
		defer fs.Close()
		WriteFile(path, ReadFile(fs))

		if ex {
			errnie.Handles(os.Chmod(path, 0777))
		}
	}
}

/*
FileExists checks if a file is present at a certain path.
*/
func FileExists(path string) bool {
	errnie.Traces()
	_, err := os.Stat(path)
	return !os.IsNotExist(err) // Have to reverse the logic.
}

/*
GetEmbedded opens the embedded file system.
*/
func GetEmbedded(embedded embed.FS, cfgFile string) fs.File {
	errnie.Traces()
	chunks := strings.Split(cfgFile, "/")

	fs, err := embedded.Open("cfg/" + chunks[len(chunks)-1])
	errnie.Handles(err)

	return fs
}

/*
ReadFile takes a file handle and reads the contents into a buffer.
*/
func ReadFile(fs fs.File) []byte {
	errnie.Traces()
	buf, err := io.ReadAll(fs)
	errnie.Handles(err)
	return buf
}

/*
WriteFile dumps a buffer to a file.
*/
func WriteFile(path string, buf []byte) {
	errnie.Traces()
	errnie.Handles(
		os.WriteFile(path, buf, 0644),
	)
}

/*
Copy a file from one location to another.
*/
func Copy(origin string, destination string) {
	errnie.Traces()
	bytesRead, err := os.ReadFile(origin)
	errnie.Handles(err)

	err = os.WriteFile(destination, bytesRead, 0755)
	errnie.Handles(err)
}

/*
DeleteFile removes a file from the filesystem.
*/
func DeleteFile(path string) {
	errnie.Traces()
	errnie.Handles(os.Remove(path))
}
