package matrix

import (
	"io"

	"github.com/theapemachine/errnie/v2"

	"github.com/docker/docker/pkg/archive"
)

type Tar struct {
	errs errnie.Collector
	src  string
}

func NewTar(src string) Tar {
	return Tar{src: src}
}

func (tar Tar) Compress() (io.ReadCloser, error) {
	return archive.TarWithOptions(tar.src, &archive.TarOptions{})
}
