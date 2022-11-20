package brazil

import (
	"os"
	"strings"

	"github.com/theapemachine/wrkspc/errnie"
)

type Path struct {
	Location string
	err      *errnie.Error
}

func NewPath(segments ...string) *Path {
	path := &Path{}
	path.Location = path.toPrefix(segments...)
	return path
}

func (path *Path) toPrefix(segments ...string) string {
	var out string
	var err error

	switch segments[0] {
	case "~":
		out, err = os.UserHomeDir()
	case ".":
		out, err = os.Getwd()
	default:
		out = strings.Join(segments, "/")
		path.makePath(out)
	}

	if errnie.Handles(err) != nil {
		path.err = errnie.NewError(err)
		return ""
	}

	return out
}

func (path *Path) makePath(prefix string) {
	_, err := os.Stat(prefix)

	if errnie.Handles(err) != nil {
		path.err = errnie.NewError(err)
		errnie.Handles(os.MkdirAll(prefix, os.ModePerm))
		errnie.Informs("created new path at", prefix)
	}

}
