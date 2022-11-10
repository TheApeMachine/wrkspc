package lumiere

import (
	"bytes"
	"strconv"

	svg "github.com/ajstarks/svgo"
)

type Button struct {
	fill  string
	state int
}

func NewButton() Element {
	return NewElement(
		Button{},
	)
}

func (button Button) Render() *bytes.Buffer {
	width := 500
	height := 250
	buffer := bytes.NewBuffer([]byte{})

	canvas := svg.New(buffer)
	canvas.Start(width, height)
	canvas.Rect(0, 0, 500, 250, button.fill)
	canvas.Text(
		width/2, height/2,
		strconv.Itoa(button.state),
		"text-anchor:middle;font-size:30px;fill:white",
	)
	canvas.End()

	return buffer
}
