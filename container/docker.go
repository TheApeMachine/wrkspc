package container

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	dc "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

type Docker struct {
	ID     string
	ctx    *twoface.Context
	cli    *client.Client
	vendor string
	name   string
	tag    string
	err    error
}

func NewDocker(
	ctx *twoface.Context, vendor, name string, tag string,
) *Docker {
	var (
		cli *client.Client
		err error
	)

	if cli, err = client.NewClientWithOpts(
		client.FromEnv, client.WithAPIVersionNegotiation(),
	); errnie.Handles(err) != nil {
		return nil
	}

	return &Docker{"", ctx, cli, vendor, name, tag, nil}
}

func (docker *Docker) Start() *Docker {
	var (
		err error
		out io.ReadCloser
	)

	if err = docker.cli.ContainerStart(
		docker.ctx, docker.ID, types.ContainerStartOptions{},
	); errnie.Handles(err) != nil {
		return docker
	}

	statusCh, errCh := docker.cli.ContainerWait(
		docker.ctx,
		docker.ID,
		dc.WaitConditionNotRunning,
	)

	select {
	case err = <-errCh:
		if errnie.Handles(err) != nil {
			return docker
		}
	case <-statusCh:
	}

	if out, err = docker.cli.ContainerLogs(
		docker.ctx, docker.ID, types.ContainerLogsOptions{
			ShowStdout: true,
		},
	); errnie.Handles(err) != nil {
		return docker
	}

	io.Copy(os.Stdout, out)
	return docker
}

func (docker *Docker) Create(entrypoint, cmd *[]string) *Docker {
	var (
		reader dc.ContainerCreateCreatedBody
		err    error
	)

	config := &dc.Config{Tty: true}

	if entrypoint != nil {
		config.Entrypoint = *entrypoint
	}

	if cmd != nil {
		config.Cmd = *cmd
	}

	if reader, err = docker.cli.ContainerCreate(
		docker.ctx.Root(),
		config,
		&dc.HostConfig{
			Binds: []string{fmt.Sprintf(
				"/tmp/wrkspc/%s/%s", docker.vendor, docker.name,
			)},
		},
		&network.NetworkingConfig{},
		&specs.Platform{},
		"",
	); errnie.Handles(err) != nil {
		return docker
	}

	errnie.Informs(fmt.Sprintf("CREATED %s", reader.ID))
	docker.ID = reader.ID
	errnie.Debugs(reader.Warnings)

	return docker
}

func (docker *Docker) Pull() *Docker {
	var (
		reader io.ReadCloser
		err    error
	)

	if reader, err = docker.cli.ImagePull(
		docker.ctx.Root(),
		fmt.Sprintf("%s/%s:%s", docker.vendor, docker.name, docker.tag),
		types.ImagePullOptions{},
	); errnie.Handles(err) != nil {
		return docker
	}

	io.Copy(os.Stdout, reader)
	return docker
}

func (docker *Docker) Build() *Docker {
	var (
		tar  io.ReadCloser
		res  types.ImageBuildResponse
		path = brazil.NewPath(
			fmt.Sprintf("/tmp/wrkspc/%s/", docker.name),
		)
	)

	if tar, docker.err = archive.TarWithOptions(
		path.Location, &archive.TarOptions{},
	); errnie.Handles(docker.err) != nil {
		return docker
	}

	if res, docker.err = docker.cli.ImageBuild(
		docker.ctx.Root(), tar, types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags: []string{fmt.Sprintf(
				"%s/%s:%s", docker.vendor, docker.name, docker.tag,
			)},
			Remove:     true,
			PullParent: true,
			Platform:   "linux/amd64",
		},
	); errnie.Handles(docker.err) != nil {
		return docker
	}

	scanner := bufio.NewScanner(res.Body)

	for scanner.Scan() {
		errnie.Debugs(scanner.Text())
	}

	return docker
}
