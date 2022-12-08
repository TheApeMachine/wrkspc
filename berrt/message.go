package berrt

type Message struct {
	From  *Object
	To    *Object
	Label string
	Note  string
}

func NewMessage(from, to *Object, label, note string) *Message {
	return &Message{from, to, label, note}
}
