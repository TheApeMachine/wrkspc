package spd

type RoleType []byte

var (
	EMPTY RoleType = []byte("empty")
	DATAP RoleType = []byte("datapoint")
	QUERY RoleType = []byte("query")
	ERROR RoleType = []byte("error")
)

type ScopeType []byte

var (
	UNMARSHAL ScopeType = []byte("unmarshal")
)
