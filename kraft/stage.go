package kraft

import "github.com/theapemachine/wrkspc/twoface"

type Stage interface {
	Make() error
}

func NewStager(ctx *twoface.Context) []Stage {
	return []Stage{
		&PkgStage{ctx},
		&BuildStage{ctx},
		&RunStage{ctx},
	}
}
