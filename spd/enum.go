package spd

type MediaType []byte
type RoleType []byte
type ScopeType []byte

var (
	APPJSN MediaType = []byte("application/json")

	EMPTY       RoleType = []byte("empty")
	DATAPOINT   RoleType = []byte("datapoint")
	QUESTION    RoleType = []byte("question")
	ERROR       RoleType = []byte("error")
	AGGREGATION RoleType = []byte("aggregation")
	CHANNEL     RoleType = []byte("channel")
	BUFFER      RoleType = []byte("buffer")

	MERGE     ScopeType = []byte("merge")
	IPC       ScopeType = []byte("ipc")
	UNMARSHAL ScopeType = []byte("unmarshal")
)
