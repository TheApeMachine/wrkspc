package fellini

import (
	"github.com/theapemachine/wrkspc/twoface"
)

/*
Console is a Template for a kubrick Layout that produces a terminal console the user can
interact with.
*/
type Console struct {
	disposer *twoface.Disposer
	input    Component
}

/*
Initialize the Console Template, setting up any state we need to operate.
*/
func (template Console) Initialize(channel chan []byte) Template {
	template.disposer = twoface.NewDisposer()
	template.input = NewInput(
		TextInput{},
		channel,
		template.disposer,
	)

	return template
}
