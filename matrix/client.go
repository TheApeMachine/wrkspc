package matrix

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
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
	Conn() *containerd.Client
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
	Disposer *twoface.Disposer
	conn     *containerd.Client
}

/*
Conn return the connection to the ContainerD socket.
TODO: This should be done in a better way to abstract away the specifics to ContainerD.
*/
func (client Containerd) Conn() *containerd.Client {
	return client.conn
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

	errnie.Handles(err).With(errnie.KILL)
	return client, image
}

/*
Push a container image to a container registry.
*/
func (client Containerd) Push(name string) {}

// /*
// Docker connects to the Docker API, which allows for all Docker cli commands to run natively
// in code, without having do do shell executions or something like that.
// */
// type Docker struct {
// 	Disposer *twoface.Disposer
// 	conn     *client.Client
// }

// /*
// Initialize constructs a Client.
// */
// func (cli Docker) Initialize() Client {
// 	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
// 	errnie.Handles(err).With(errnie.KILL)
// 	cli.conn = c
// 	return cli
// }

// /*
// Pull a container image.
// */
// func (cli Docker) Pull(name string) {
// }

// /*
// Push a container image to a registry.
// */
// func (cli Docker) Push(name string) {
// 	creds := struct {
// 		Username      string `json:"username"`
// 		Password      string `json:"password"`
// 		Email         string `json:"email"`
// 		Serveraddress string `json:"serveraddress"`
// 	}{
// 		Username:      viper.GetString("wrkspc.matrix.registry.username"),
// 		Password:      viper.GetString("wrkspc.matrix.registry.password"),
// 		Email:         viper.GetString("wrkspc.matrix.registry.email"),
// 		Serveraddress: viper.GetString("wrkspc.matrix.registry.host"),
// 	}

// 	d, err := json.Marshal(creds)
// 	errnie.Handles(err).With(errnie.KILL)
// 	sEnc := base64.StdEncoding.EncodeToString([]byte(d))

// 	opts := types.ImagePushOptions{}
// 	opts.RegistryAuth = sEnc

// 	_, err = cli.conn.ImagePush(cli.Disposer.Ctx, name, opts)
// 	errnie.Handles(err).With(errnie.KILL)
// }
