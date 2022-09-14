package sockpuppet

type Conn interface {
	Up(string) error
}

func NewConn(connType Conn) Conn {
	return connType
}
