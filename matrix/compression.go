package matrix

import (
	"io"

	"github.com/docker/docker/pkg/archive"
	"github.com/theapemachine/wrkspc/errnie"
)

type Tar struct {
	src string
}

func NewTar(src string) Tar {
	errnie.Traces()
	return Tar{src: src}
}

func (tar Tar) Compress() (io.ReadCloser, error) {
	errnie.Traces()
	return archive.TarWithOptions(tar.src, &archive.TarOptions{})
}
