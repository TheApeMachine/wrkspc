package nomad

import "github.com/theapemachine/wrkspc/infra"

type Client struct{}

func NewClient() infra.Client {
	return infra.NewClient(Client{})
}

func (client Client) Apply(name, vendor, namespace string) {
}
