package tui

/*
Ansi is an object that improves developer ergonomics around executing
Ansi escape requences for advanced control over the terminal.
*/
type Ansi struct{}

func NewAnsi() *Ansi {
	return &Ansi{}
}

func (screen *Ansi) ToggleAltScreen()
