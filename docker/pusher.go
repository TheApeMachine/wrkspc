package docker

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Pusher is a type responsible for pushing container images
to a registry so that they can be accessed from any machine
that has access to the registry.
*/
type Pusher struct {
	vendor string
	name   string
	tag    string
	ctx    *twoface.Context
	cli    *client.Client
}

/*
NewPusher returns a pointer to a Pusher type.
*/
func NewPusher(vendor, name, tag string) *Pusher {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	errnie.Handles(err)

	return &Pusher{
		vendor: vendor,
		name:   name,
		tag:    tag,
		ctx:    twoface.NewContext(),
		cli:    cli,
	}
}

/*
Send the Docker image to a registry.
This makes sure that no matter where you are or on which machine,
you will always receive your current container state.
*/
func (pusher *Pusher) Send() {
	authConfigBytes, err := json.Marshal(types.AuthConfig{
		Username:      viper.GetString("docker.username"),
		Password:      viper.GetString("docker.password"),
		ServerAddress: "https://index.docker.io/v1/",
	})

	errnie.Handles(err)

	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)

	rd, err := pusher.cli.ImagePush(
		pusher.ctx.Handle(),
		pusher.vendor+"/"+pusher.name+":"+pusher.tag,
		types.ImagePushOptions{RegistryAuth: authConfigEncoded},
	)

	errnie.Handles(err)
	defer rd.Close()

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	errnie.Handles(scanner.Err())
}
