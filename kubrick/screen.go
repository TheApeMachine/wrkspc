package kubrick

/*
Screen is a container object for a collection of Layouts and is responsible for high level operations
regarding the display of the current Layout's content, switching layouts, etc.
*/
type Screen struct {
	layouts []Layout
}

/*
NewScreen constructs a new Screen using the Layouts that are passed in.
*/
func NewScreen(layouts ...Layout) *Screen {
	return &Screen{layouts: layouts}
}

/*
Render the current Layout.
*/
func (screen *Screen) Render() chan error {
	out := make(chan error)

	go func() {
		defer close(out)
		out <- <-screen.layouts[0].Render()
	}()

	return out
}
