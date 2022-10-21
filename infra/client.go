package infra

/*
Client represents the connection to a cluster.
*/
type Client interface {
	Apply(string, string, string)
}

/*
NewClient converts a cluster client struct type to its interface
representation.
*/
func NewClient(clientType Client) Client {
	return clientType
}
