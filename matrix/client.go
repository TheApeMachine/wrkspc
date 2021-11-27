package matrix

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Client is an interface to be implemented by objects that want to be a client to a container API.
*/
type Client interface {
	Pull(string) (Client, containerd.Image)
	Push(string)
	ToSpec(context.Context, string, containerd.Image) *specs.Spec
	Conn() *containerd.Client
	Cleanup(context.Context)
}

/*
NewClient constructs a Client of the type that is passed in.
*/
func NewClient(clientType Client) Client {
	return clientType
}

/*
Containerd is a connection to the ContainerD daemon that is running in the background.
*/
type Containerd struct {
	Disposer  *twoface.Disposer
	conn      *containerd.Client
	container containerd.Container
}

/*
Conn return the connection to the ContainerD socket.
TODO: This should be done in a better way to abstract away the specifics to ContainerD.
*/
func (client Containerd) Conn() *containerd.Client {
	errnie.Traces()
	return client.conn
}

/*
ToSpec generates an OCI spec from the container image.
*/
func (client Containerd) ToSpec(
	ctx context.Context, name string, image containerd.Image,
) *specs.Spec {
	errnie.Traces()

	searchSpace, err := client.conn.Containers(ctx, name)

	errnie.Handles(err).With(errnie.NOOP)
	errnie.Logs(searchSpace).With(errnie.DEBUG)

	if len(searchSpace) == 0 {
		errnie.Logs("making new container").With(errnie.INFO)
		client.container, err = client.conn.NewContainer(
			ctx, name,
			containerd.WithNewSnapshot(name+"-snapshot1", image),
			containerd.WithNewSpec(oci.WithImageConfig(image)),
		)
		errnie.Handles(err).With(errnie.NOOP)
		errnie.Logs(client.container).With(errnie.DEBUG)
	} else {
		client.container = searchSpace[0]
	}

	spec, err := client.container.Spec(ctx)
	errnie.Handles(err).With(errnie.NOOP)
	errnie.Traces()

	return spec
}

/*
Cleanup the Container.
*/
func (client Containerd) Cleanup(ctx context.Context) {
	errnie.Traces()
	defer client.conn.Close()
	client.container.Delete(ctx, containerd.WithSnapshotCleanup)
}

/*
Pull a container image from a container registry.
*/
func (client Containerd) Pull(name string) (Client, containerd.Image) {
	var err error

	if client.conn == nil {
		client.conn, err = containerd.New("/run/containerd/containerd.sock")
		errnie.Handles(err).With(errnie.KILL)
		errnie.Logs("connected to containerd socket ", client.conn).With(errnie.INFO)
	}

	// We need to upgrade the context with a namespace for the ContainerD daemon to work.
	client.Disposer.Ctx = namespaces.WithNamespace(context.Background(), "test1")

	image, err := client.conn.Pull(
		client.Disposer.Ctx,
		viper.GetString("wrkspc.matrix.registry.host")+"/"+
			viper.GetString("wrkspc.matrix.registry.username")+"/"+name+":v1.0",
		// TODO: remove harcoded latest tag, just want containers to run again first.
		containerd.WithPullUnpack,
	)

	errnie.Handles(err).With(errnie.NOOP)
	return client, image
}

/*
Push a container image to a container registry.
*/
func (client Containerd) Push(name string) {}
