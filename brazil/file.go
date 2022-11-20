package brazil

import (
	"bytes"
	"os"
	"strings"

	"github.com/theapemachine/wrkspc/errnie"
)

type File struct {
	Location string
	Name     string
	Ext      string
	Data     *bytes.Buffer
	err      *errnie.Error
}

func NewFile(path, name string, data *bytes.Buffer) *File {
	errnie.Trace()

	nSplit := strings.Split(name, ".")
	ext := nSplit[len(nSplit)-1]
	fName := name[:len(name)-len(nSplit)-1]

	fh := &File{
		Location: path,
		Name:     name,
		Ext:      ext,
		Data:     data,
	}

	fullPath := strings.Join(
		[]string{strings.Join([]string{path, fName}, "/"), ext}, ".",
	)

	fh.createIfNotExists(fullPath, data)
	buf, err := os.ReadFile(fullPath)

	if fh.err = errnie.Handles(err); fh.err != nil {
		return nil
	}

	fh.Data = bytes.NewBuffer(buf)
	return fh
}

func (file *File) createIfNotExists(fullPath string, data *bytes.Buffer) {
	errnie.Trace()

	_, err := os.Stat(fullPath)

	if os.IsNotExist(err) {
		file.err = errnie.Handles(errnie.NewError(err))
		errnie.Handles(os.WriteFile(fullPath, data.Bytes(), 0644))
		errnie.Informs("copied config to", fullPath)
	}
}
