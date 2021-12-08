package matrix

import (
	"context"
	"os"

	"github.com/containerd/containerd"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/client/llb/imagemetaresolver"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Build represents a container that is built.
*/
type Build struct {
	name   string
	tag    string
	client Client
	root   string
	image  containerd.Image
}

/*
NewBuild kicks off a new container build, pass in the name you want to give the container.
*/
func NewBuild(name, tag string, client Client) *Build {
	return &Build{
		name:   name,
		tag:    tag,
		client: client,
		root: brazil.BuildPath(
			brazil.HomePath(), ".wrkspc", "manifests", "dockerfiles",
		),
	}
}

/*
SetImage sets or overrides the image of the build, this is mostly used for manual builds.
*/
func (build *Build) SetImage(image containerd.Image) *Build {
	build.image = image
	return build
}

/*
ToLLB converts a Dockerfile to an LLB spec so BuildKit has something to work with that it likes.
*/
func (build *Build) ToLLB(name, tag string) *dockerfile2llb.Image {
	errnie.Traces()

	dockerfile := brazil.NewFile(brazil.BuildPath(
		brazil.HomePath(), "wrkspc", name, "Dockerfile",
	))

	caps := pb.Caps.CapSet(pb.Caps.All())

	state, image, err := dockerfile2llb.Dockerfile2LLB(
		appcontext.Context(), dockerfile.Data.Bytes(), dockerfile2llb.ConvertOpt{
			MetaResolver: imagemetaresolver.Default(),
			Target:       name + ":" + tag,
			LLBCaps:      &caps,
		},
	)

	errnie.Handles(err).With(errnie.KILL)

	dt, err := state.Marshal(context.TODO())
	errnie.Handles(err).With(errnie.KILL)

	llb.WriteTo(dt, os.Stdout)

	return image
}
