package docker

import (
	"context"
	"os"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/client/llb/imagemetaresolver"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
)

type Builder struct {
	org  string
	name string
	tag  string
}

func NewBuilder(org, name, tag string) *Builder {
	return &Builder{
		org:  org,
		name: name,
		tag:  tag,
	}
}

/*
ToLLB converts a Dockerfile to an image format that is compatible
with BuildKit, so we can have parallel build stages and be faster.
*/
func (builder *Builder) ToLLB(name, tag string) *dockerfile2llb.Image {
	dockerfile := brazil.NewFile(brazil.Workdir())
	caps := pb.Caps.CapSet(pb.Caps.All())

	state, image, _, err := dockerfile2llb.Dockerfile2LLB(
		appcontext.Context(), dockerfile.Data.Bytes(), dockerfile2llb.ConvertOpt{
			MetaResolver: imagemetaresolver.Default(),
			Target:       name + ":" + tag,
			LLBCaps:      &caps,
			BuildArgs: map[string]string{
				"USERNAME": os.Getenv("USER"),
			},
		},
	)

	errnie.Handles(err)
	dt, err := state.Marshal(context.TODO())
	errnie.Handles(err)

	llb.WriteTo(dt, os.Stdout)

	return image
}
