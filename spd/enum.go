package spd

type MediaType []byte
type RoleType []byte
type ScopeType []byte

var (
	APPXMP MediaType = []byte("application/example")
	APPBIN MediaType = []byte("application/octet-stream")
	APPTXT MediaType = []byte("application/text")
	APPJSN MediaType = []byte("application/json")

	EMPTY     RoleType = []byte("empty")
	ERROR     RoleType = []byte("error")
	SECURITY  RoleType = []byte("error")
	TEST      RoleType = []byte("test")
	DATAPOINT RoleType = []byte("datapoint")
	QUESTION  RoleType = []byte("question")
	REQUEST   RoleType = []byte("request")
	PIPE      RoleType = []byte("pipe")
	LINK      RoleType = []byte("link")

	UNKNOWN    ScopeType = []byte("unknown")
	IO         ScopeType = []byte("io")
	VALIDATION ScopeType = []byte("validation")
	ORIGIN     ScopeType = []byte("origin")
	UNIT       ScopeType = []byte("unit")
	BENCHMARK  ScopeType = []byte("benchmark")
	WORKSPACE  ScopeType = []byte("workspace")
	DATALAKE   ScopeType = []byte("datalake")
	HTTP       ScopeType = []byte("http")
	WAN        ScopeType = []byte("wan")
	UI         ScopeType = []byte("ui")
	ERRNIE     ScopeType = []byte("errnie")
)
