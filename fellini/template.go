package fellini

/*
Template is an interface that objects can implement so they become usable by a kubrick Layout.
A template defines a predefined set of Components which build up to some terminal UI functionality.
*/
type Template interface {
	Initialize(chan []byte) Template
}

/*
NewTemplate constructs a new Template of the type passed in and calls its Initialize method where
any setup can be performed.
*/
func NewTemplate(templateType Template, channel chan []byte) Template {
	return templateType.Initialize(channel)
}
