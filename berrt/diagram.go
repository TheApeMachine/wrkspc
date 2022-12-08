package berrt

import "strings"

/*
Diagram is an interface that all objects within `berrt` have to
implement, so they can be composed together.
*/
type Diagram interface {
	Render() string
}

type SequenceDiagram struct {
	lines []Diagram
}

func NewSequenceDiagram() *SequenceDiagram {
	return &SequenceDiagram{}
}

func (diagram *SequenceDiagram) Render() string {
	var builder strings.Builder
	builder.WriteString("@startuml\n")

	for _, line := range diagram.lines {
		builder.WriteString(line.Render())
	}

	builder.WriteString("@enduml\n")
	return builder.String()
}
