package spd

type DataType []byte
type RoleType []byte
type ScopeType []byte

var (
	// DataType enumerates MIME types.
	APPJSON  DataType = []byte("application/json")
	PLAINTXT DataType = []byte("text/plain")
	TXTHTML  DataType = []byte("text/html")
)

var (
	// RoleType enumerates the roles of the nodes in the network.
	DATA     RoleType = []byte("data")
	ARTIFACT RoleType = []byte("artifact")
)

var (
	// ScopeType enumerates the scope of the data.
	SCRAPE ScopeType = []byte("scrape")
)
