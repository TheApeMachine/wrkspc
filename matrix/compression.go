package matrix

import (
	"io"

	"github.com/docker/docker/pkg/archive"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Tar is a wrapper around tar based compression.
*/
type Tar struct {
	src string
}

/*
NewTar wraps a new tarball.
*/
func NewTar(src string) Tar {
	errnie.Traces()
	return Tar{src: src}
}

/*
Compress data into the tarball.
*/
func (tar Tar) Compress() (io.ReadCloser, error) {
	errnie.Traces()
	return archive.TarWithOptions(tar.src, &archive.TarOptions{})
}
