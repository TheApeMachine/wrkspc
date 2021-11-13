package matrix

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Client connects to the Docker API, which allows for all Docker cli commands to run natively
in code, without having do do shell executions or something like that.
*/
type Client struct {
	ctx  context.Context
	conn *client.Client
}

/*
NewClient constructs a Client.
*/
func NewClient(ctx context.Context) Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	errnie.Handles(err).With(errnie.KILL)

	return Client{
		ctx:  ctx,
		conn: cli,
	}
}

/*
Push a container image to a registry.
*/
func (client Client) Push(name string) {
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

	_, err = client.conn.ImagePush(client.ctx, name, opts)
	errnie.Handles(err).With(errnie.KILL)
}
