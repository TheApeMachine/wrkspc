package sherlock

/*
Lead is a collection of clues and/or connections that
need to be investigated.
*/
type Lead struct {
}

/*
NewLead creates a new Lead object.
*/
func NewLead() *Lead {
	return &Lead{}
}