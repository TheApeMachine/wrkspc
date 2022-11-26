package spd

type RoleType []byte

var (
	ERROR RoleType = []byte("error")
)

type ScopeType []byte

var (
	UNMARSHAL ScopeType = []byte("unmarshal")
)
