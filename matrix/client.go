package matrix

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
	"github.com/theapemachine/wrkspc/valleygirl"
)

/*
Client is an interface to be implemented by objects that want to be a client to a container API.
*/
type Client interface {
	Dial() Client
	Cleanup()
	Fetch(string, string) (containerd.Container, *specs.Spec)
	ToSpec(string, containerd.Image) (containerd.Container, *specs.Spec)
}

/*
NewClient constructs a Client of the type that is passed in.
*/
func NewClient(clientType Client) Client {
	errnie.Traces()
	return clientType.Dial()
}

/*
Containerd is a connection to the ContainerD daemon that is running in the background.
*/
type Containerd struct {
	Disposer *twoface.Disposer
	conn     *containerd.Client
}

/*
Dial the Client into an endpoint.
*/
func (client Containerd) Dial() Client {
	errnie.Traces()

	if client.conn == nil {
		// Since removing the hard 3 second delay in the daemon boot process we need
		// a way to reliably retry the connection until Containerd is booted fully.
		// We can use a repeater for that using our favorite backoff strategy.
		repeater := twoface.NewRepeater(10, twoface.Fibonacci{MaxTries: 10})
		repeater.Attempt(1, func() bool {
			var err error
			client.conn, err = containerd.New("/run/containerd/containerd.sock")
			return errnie.Handles(err).With(errnie.NOOP).OK
		}) // This will loop max 10 times, each time with a timeout of the two previous
		// timeout values summed together.

		// We need to upgrade the context with a namespace for the ContainerD daemon to work.
		client.Disposer.Ctx = namespaces.WithNamespace(context.Background(), "wrkspc")
	}

	return client
}

/*
Cleanup the Container.
*/
func (client Containerd) Cleanup() {
	errnie.Traces()
	defer client.conn.Close()
}

/*
Fetch a runnable container by either building it from local manifest or pulling it from
a remote registry configured in `~/.wrkspc.yml`.
*/
func (client Containerd) Fetch(name, tag string) (containerd.Container, *specs.Spec) {
	errnie.Traces()
	// var llbimage *dockerfile2llb.Image
	var image containerd.Image

	if brazil.FileExists(brazil.BuildPath(
		brazil.HomePath(), ".wrkspc", "manifests", "dockerfiles", name, "Dockerfile",
	)) {
		// builder := NewBuild(name, tag, client)
		// llbimage = builder.ToLLB(name, tag)
	} else {
		image = client.pull(name, tag)
	}

	return client.ToSpec(name, image)
}

/*
pull a container image from a container registry.
*/
func (client Containerd) pull(name, tag string) containerd.Image {
	errnie.Traces()

	image, err := client.conn.Pull(
		client.Disposer.Ctx,
		brazil.BuildPath(
			viper.GetString("wrkspc.matrix.registry.host"),
			viper.GetString("wrkspc.matrix.registry.username"),
			name+":"+tag,
		),
		containerd.WithPullUnpack,
	)

	errnie.Handles(err).With(errnie.NOOP)
	errnie.Logs("pulled image", image)

	return image
}

/*
ToSpec generates an OCI spec from the container image.
*/
func (client Containerd) ToSpec(
	name string, image containerd.Image,
) (containerd.Container, *specs.Spec) {
	errnie.Traces()

	container, err := client.conn.NewContainer(
		client.Disposer.Ctx, name,
		containerd.WithNewSnapshot(
			name+"-"+valleygirl.FastRandomString(6), image,
		),
		containerd.WithNewSpec(oci.WithImageConfig(image), oci.WithTTY),
	)

	errnie.Handles(err).With(errnie.NOOP)

	spec, err := container.Spec(client.Disposer.Ctx)
	errnie.Handles(err).With(errnie.NOOP)

	return container, spec
}
