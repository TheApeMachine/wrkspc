package matrix

import (
	"io"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/docker/docker/api/types"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Image is a respresentation of a Docker Image and can be
built to run as a container.
*/
type Image struct {
	disposer *twoface.Disposer
	name     string
	pkg      io.Reader
	opts     types.ImageBuildOptions
}

/*
NewImage returns an instance of image that is ready to be built.
It receives a tarball package that was generated from a build
context (a directory with a Dockerfile and other dependencies).
*/
func NewImage(disposer *twoface.Disposer, name string, pkg io.Reader) Image {
	errnie.Traces()

	root := viper.GetString("wrkspc.matrix.registry.host")
	username := viper.GetString("wrkspc.matrix.registry.username")

	imgTag := []string{root + "/" + username + "/" + name}

	return Image{
		disposer: disposer,
		name:     name,
		pkg:      pkg,
		opts: types.ImageBuildOptions{
			Context: pkg,
			Tags:    imgTag,
			Remove:  true,
			BuildArgs: map[string]*string{
				"USER": &username,
			},
		},
	}
}

/*
Build the final image and return the log stream from the
builder daemon.
*/
func (img Image) Build(cli Client) containerd.Container {
	ctx := namespaces.WithNamespace(img.disposer.Ctx, img.name)

	image := cli.Pull(img.name)

	build, err := cli.Conn().NewContainer(
		ctx, img.name,
		containerd.WithNewSnapshot(img.name+"-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)

	errnie.Handles(err).With(errnie.KILL)
	return build
}
