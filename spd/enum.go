package spd

type RoleType []byte

var (
	EMPTY       RoleType = []byte("empty")
	DATAPOINT   RoleType = []byte("datapoint")
	QUESTION    RoleType = []byte("question")
	ERROR       RoleType = []byte("error")
	AGGREGATION RoleType = []byte("aggregation")
	CHANNEL     RoleType = []byte("channel")
)

type ScopeType []byte

var (
	IPC       ScopeType = []byte("ipc")
	UNMARSHAL ScopeType = []byte("unmarshal")
)
