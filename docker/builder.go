package docker

import (
	"bufio"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
	"github.com/moby/moby/pkg/archive"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type Builder struct {
	org  string
	name string
	tag  string
	ctx  *twoface.Context
	cli  *client.Client
}

func NewBuilder(org, name, tag string) *Builder {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	errnie.Handles(err)

	return &Builder{
		org:  org,
		name: name,
		tag:  tag,
		ctx:  twoface.NewContext(),
		cli:  cli,
	}
}

func (builder *Builder) Build() {
	tar, err := archive.TarWithOptions(
		brazil.BuildPath(brazil.Workdir()),
		&archive.TarOptions{},
	)

	errnie.Handles(err)

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{builder.org, builder.name},
		Remove:     true,
	}

	res, err := builder.cli.ImageBuild(builder.ctx.Handle(), tar, opts)
	errnie.Handles(err)
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	errnie.Handles(scanner.Err())
}

/*
ToLLB converts a Dockerfile to an image format that is compatible
with BuildKit, so we can have parallel build stages and be faster.
*/
/*
func (builder *Builder) ToLLB() *dockerfile2llb.Image {
	dockerfile := brazil.NewFile(brazil.Workdir())
	caps := pb.Caps.CapSet(pb.Caps.All())

	target := strings.Builder{}
	target.WriteString(builder.org)
	target.WriteString("/")
	target.WriteString(builder.name)
	target.WriteString(":")
	target.WriteString(builder.tag)

	state, image, err := dockerfile2llb.Dockerfile2LLB(
		appcontext.Context(), dockerfile.Data.Bytes(), dockerfile2llb.ConvertOpt{
			MetaResolver: imagemetaresolver.Default(),
			Target:       target.String(),
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
*/
