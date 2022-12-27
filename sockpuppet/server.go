package sockpuppet

/*
Server is an interface for objects to implement that want
to accept Connections from some Client.
*/
type Server interface {
	Up(string) error
}
