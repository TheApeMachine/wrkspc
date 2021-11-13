package matrix

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Image is a respresentation of a Docker Image and can be
built to run as a container.
*/
type Image struct {
	ctx  context.Context
	name string
	pkg  io.Reader
	opts types.ImageBuildOptions
}

/*
NewImage returns an instance of image that is ready to be built.
It receives a tarball package that was generated from a build
context (a directory with a Dockerfile and other dependencies).
*/
func NewImage(ctx context.Context, name string, pkg io.Reader, outname string) Image {
	root := viper.GetString("wrkspc.matrix.registry.host")
	username := viper.GetString("wrkspc.matrix.registry.username")

	imgTag := []string{root + "/" + username + "/" + outname}

	return Image{
		ctx:  ctx,
		name: name,
		pkg:  pkg,
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
func (img Image) Build(cli Client) types.ImageBuildResponse {
	res, err := cli.conn.ImageBuild(img.ctx, img.pkg, img.opts)
	errnie.Handles(err).With(errnie.KILL)

	return res
}
