package matrix

import (
	"github.com/containerd/containerd"
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
Atomic defines a high-level flow of build steps that each build exactly the same
way every time, or fail. A failure usually means somebody has done something
manually somewhere.
*/
func (build *Build) Atomic(fs bool) error {
	errnie.Traces()

	// Since we are always rebuilding the root filesystem for the image, we need
	// to have a way to skip it during that stage to prevent an infinite loop.
	if !fs {
		// Make sure we use the current base image for the tool we are building.
		rootfs := NewRootFS(build)
		rootfs.Build()

		// Update all the dependencies in our ~/.wrkspc path
		resolver := NewResolver(build.root)
		resolver.Update()
	}

	// Wrap the context up into a tarball to send to the builder daemon.
	tar := NewTar(brazil.BuildPath(build.root, build.name))

	pkg, err := tar.Compress()
	_ = pkg
	errnie.Handles(err).With(errnie.KILL)

	// // Prepare a new image spec for thedaemon to build.
	// build.image = NewImage(build.name, pkg)

	// // TODO: This is changing a lot due to ContainerD vs Docker integration. Will clean up later.
	// // Tell the daemon to build our image.
	// build.container = build.image.Build(build.client)
	// // scanner := NewScanner(build.img)
	// // scanner.Scan()

	// errnie.Logs("container", build.container).With(errnie.DEBUG)

	// cspec, err := build.container.Spec(build.disposer.Ctx)
	// errnie.Handles(err).With(errnie.NOOP)
	// errnie.Logs("spec", cspec).With(errnie.DEBUG)

	// containerProto, err := build.container.Info(build.disposer.Ctx)
	// errnie.Handles(err).With(errnie.NOOP)

	// errnie.Handles(oci.WithTTY(
	// 	build.disposer.Ctx, client.Conn(), &containerProto, cspec,
	// ))

	// run := NewRun(build, cspec)
	// run.Cycle(build.disposer.Ctx)

	return err
}
