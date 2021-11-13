package spdg

/*
Annotation is a simple key/value pair that acts as a parameter to
the Role of the Datagram.
*/
type Annotation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

/*
NewAnnotation constructs an Annoation.
*/
func NewAnnotation(key, value string) Annotation {
	return Annotation{key, value}
}
