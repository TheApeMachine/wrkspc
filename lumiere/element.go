package lumiere

import "bytes"

/*
Element is an interface for widgets to implement.
*/
type Element interface {
	Render() *bytes.Buffer
}

/*
NewElement converts a widget structure and returns its interface type.
*/
func NewElement(elementType Element) Element {
	return elementType
}
