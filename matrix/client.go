package matrix

import (
	"encoding/base64"
	"encoding/json"

	"github.com/containerd/containerd"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Client is an interface to be implemented by objects that want to be a client to a container API.
*/
type Client interface {
	Pull(string)
	Push(string)
}

/*
NewClient constructs a Client of the type that is passed in.
*/
func NewClient(clientType Client) Client {
	return clientType
}

type Containerd struct {
	Disposer *twoface.Disposer
}

func (client Containerd) Pull(name string) {
	conn, err := containerd.New("/run/containerd/containerd.sock")
	errnie.Handles(err).With(errnie.KILL)
	defer conn.Close()

	image, err := client.Pull(
		client.Disposer.Ctx,
		viper.GetString("wrkspc.matrix.registry.host/")+
			viper.GetString("wrkspc.matrix.registry.username/")+name,
		containerd.WithPullUnpack,
	)

	errnie.Handles(err).With(errnie.KILL)
	return image
}

/*
Docker connects to the Docker API, which allows for all Docker cli commands to run natively
in code, without having do do shell executions or something like that.
*/
type Docker struct {
	Disposer *twoface.Disposer
	conn     *client.Client
}

/*
Initialize constructs a Client.
*/
func (cli Docker) Initialize() Client {
	c, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	errnie.Handles(err).With(errnie.KILL)
	cli.conn = c
	return cli
}

/*
Pull a container image.
*/
func (cli Docker) Pull(name string) {
}

/*
Push a container image to a registry.
*/
func (cli Docker) Push(name string) {
	creds := struct {
		Username      string `json:"username"`
		Password      string `json:"password"`
		Email         string `json:"email"`
		Serveraddress string `json:"serveraddress"`
	}{
		Username:      viper.GetString("wrkspc.matrix.registry.username"),
		Password:      viper.GetString("wrkspc.matrix.registry.password"),
		Email:         viper.GetString("wrkspc.matrix.registry.email"),
		Serveraddress: viper.GetString("wrkspc.matrix.registry.host"),
	}

	d, err := json.Marshal(creds)
	errnie.Handles(err).With(errnie.KILL)
	sEnc := base64.StdEncoding.EncodeToString([]byte(d))

	opts := types.ImagePushOptions{}
	opts.RegistryAuth = sEnc

	_, err = cli.conn.ImagePush(cli.Disposer.Ctx, name, opts)
	errnie.Handles(err).With(errnie.KILL)
}
