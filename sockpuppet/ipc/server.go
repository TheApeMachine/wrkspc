package ipc

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (server *Server) Up(port string) error {
	return nil
}
