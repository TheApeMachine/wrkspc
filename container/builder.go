package container

type Builder struct {
	Options *BuildOptions
}

type BuildOptions struct {
	Vendor string
	Name   string
	Tag    string
	Cmd    string
}

func NewBuilder(options *BuildOptions) *Builder {
	return &Builder{
		Options: options,
	}
}
