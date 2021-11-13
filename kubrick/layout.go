package kubrick

import "github.com/theapemachine/wrkspc/fellini"

/*
Layout defines the interface an object can implement to become a type that positions Components
onto a Screen.
*/
type Layout interface {
	Render() chan error
}

/*
NewLayout constructs a layout of the type that is passed in.
*/
func NewLayout(layoutType Layout) Layout {
	return layoutType
}

/*
FullScreenLayout contains no Grid or other positional helper elements and is meant to use all
available Screen realestate for a single graphical context.
*/
type FullScreenLayout struct {
	Template fellini.Template
}

/*
Render the content of a Layout to the Screen.
*/
func (layout FullScreenLayout) Render() chan error {
	out := make(chan error)
	return out
}
