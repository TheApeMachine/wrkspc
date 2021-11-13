package contempt

/*
Connector is an interface for objects to implement that know a certain Connection protocol
and want to try to connect to all found network entities using the Scanner. They will often
have multiple ways to try and establish a Connection.
*/
type Connector interface {
	Sweep() Connection
}

/*
NewConnector constructs a Connector using the type that is passed in.
*/
func NewConnector(connectorType Connector) Connector {
	return connectorType
}
