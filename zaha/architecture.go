package zaha

import (
	"github.com/theapemachine/wrkspc/passepartout"
	"github.com/theapemachine/wrkspc/sockpuppet"
)

type Service uint

const (
	GATEWAY Service = iota
)

var services = map[Service]sockpuppet.Conn{
	GATEWAY: sockpuppet.NewHTTP(passepartout.NewRouter()),
}

/*
Architecture is a generic structure with which to define a service.
It needs some form of connection object, some data store, and one
or more jobs to be able to provide a network endpoint that moves
data upon request.
*/
type Architecture struct {
	service Service
}

/*
NewArchitecture constructs a service. This should be done using a
self-contained CLI command in the ./cmd folder so that this project
can be specifically run as that service using a container.
*/
func NewArchitecture(service Service) *Architecture {
	return &Architecture{service: service}
}

/*
Build the architecture so that it is ready to be served using the serve
cli command. This allows the binary to be deployed as many services.
*/
func (architecture *Architecture) Build() sockpuppet.Conn {
	return services[architecture.service]
}
