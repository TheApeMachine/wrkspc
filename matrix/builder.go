package matrix

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Build represents a container that is built.
*/
type Build struct {
	disposer *twoface.Disposer
	root     string
	name     string
	img      types.ImageBuildResponse
}

/*
NewBuild kicks off a new container build, pass in the name you want to give the container.
*/
func NewBuild(name string) *Build {
	return &Build{
		name:     name,
		disposer: twoface.NewDisposer(),
		root:     viper.GetString("homepath") + "/.wrkspc",
	}
}

/*
Validate the Build by starting a pull from a registry or finding the command locally.
*/
func (build *Build) Validate() chan string {
	errnie.Traces()
	out := make(chan string)
	return out
}

/*
Atomic defines a high-level flow of build steps that each build exactly the same
way every time, or fail. A failure usually means somebody has done something
manually somewhere.
*/
func (build *Build) Atomic(fs bool) error {
	if viper.GetBool("wrkspc.errnie.debug") {
		wd, err := os.Getwd()
		errnie.Handles(err).With(errnie.KILL)
		build.root = wd + "/dockerfiles"
	}

	infix := ""

	if build.root != "" {
		infix = "/"
	}

	// Since we are always rebuilding the root filesystem for the image, we need
	// to have a way to skip it during that stage to prevent an infinite loop.
	if !fs {
		// Make sure we use the current base image for the tool we are building.
		rootfs := NewRootFS(build.root, build.name)
		rootfs.Build()

		// Update all the dependencies in our ~/.wrkspc path
		resolver := NewResolver(build.root)
		resolver.Update()
	}

	// Wrap the context up into a tarball to send to the builder daemon.
	outpath := filepath.FromSlash(build.root + infix + build.name + "/")
	tar := NewTar(outpath)

	pkg, err := tar.Compress()
	errnie.Handles(err).With(errnie.KILL)

	// Prepare a new image spec for the daemon to build.
	spec := NewImage(build.disposer.Ctx, build.name, pkg, build.name)

	// Get a client to the daemon so we can send our spec.
	client := NewClient(Docker{
		Disposer: build.disposer,
	})

	// Tell the daemon to build our image.
	build.img = spec.Build(client)
	scanner := NewScanner(build.img)
	scanner.Scan()

	return err
}
