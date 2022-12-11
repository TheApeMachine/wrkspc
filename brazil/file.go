package brazil

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/theapemachine/wrkspc/errnie"
)

/*
File wraps the os.File struct and provides a specialized workflow.

When referencing a file in wrkspc it will either be retrieved,
or created on-the-fly if it does not exist yet.
*/
type File struct {
	Location string
	Name     string
	Ext      string
	Info     os.FileInfo
	Data     *bytes.Buffer
	err      error
}

/*
NewFile is a constructor returning a pointer to a File instance.

When using this constructor you can be sure that everything is
handled correctly any time you reference a file.
*/
func NewFile(path, name string, data *bytes.Buffer) *File {
	errnie.Trace()

	// Split the extension off the file name, but make sure
	// to retain all other elements in fName, which can
	// include dots (.) for private files.
	nSplit := strings.Split(name, ".")
	ext := nSplit[len(nSplit)-1]
	fName := name[:len(name)-len(ext)-1]

	fh := &File{
		Location: path,
		Name:     name,
		Ext:      ext,
		Data:     data,
	}

	// Get a fully qualified path to the file handle, and create the file
	// if it does not already exist. Then read the context into a buffer.
	fullPath := strings.Join([]string{strings.Join([]string{path, fName}, "/"), ext}, ".")
	fh.createIfNotExists(fullPath, data)
	buf, err := os.ReadFile(fullPath)

	if fh.err = errnie.Handles(err); fh.err != nil {
		return nil
	}

	// Add the file data to our object and return it.
	fh.Data = bytes.NewBuffer(buf)
	return fh
}

/*
createIfNotExists writes a new file to the file system, if it does not
already exist, using the buffer we pass in to fill the new file with data.
*/
func (file *File) createIfNotExists(fullPath string, data *bytes.Buffer) {
	errnie.Trace()

	if file.Info, file.err = os.Stat(fullPath); file.err != nil {
		errnie.Warns(fmt.Sprintf("%s does not exist, creating...", fullPath))
		path := strings.Split(fullPath, "/")
		NewPath(path[:len(path)-1]...)
		errnie.Handles(os.WriteFile(fullPath, data.Bytes(), 0644))
		errnie.Informs("write file", fullPath)
	}
}
