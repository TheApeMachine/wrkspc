package berrt

type ContextType string

const (
	ALTELSE  ContextType = "alt"
	OPT      ContextType = "opt"
	LOOP     ContextType = "loop"
	PAR      ContextType = "par"
	BREAK    ContextType = "break"
	CRITICAL ContextType = "critical"
	GROUP    ContextType = "group"
)

type Context struct {
	Name     string
	Type     ContextType
	Messages []*Message
}

func NewContext(
	name string, Type ContextType, messages []*Message,
) *Context {
	return &Context{name, Type, messages}
}
