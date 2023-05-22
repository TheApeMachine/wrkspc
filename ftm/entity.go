package ftm

/*
Entity is a wrapper that mirrors the FollowTheMoney Entity object.
*/
type Entity struct {
	ID         string
	Schema     string
	Properties map[string][]string
}

/*
NewEntity creates a new Entity object.
*/
func NewEntity() *Entity {
	return &Entity{}
}

/*
Read implements the io.Reader interface.
*/
func (e *Entity) Read(p []byte) (n int, err error) {
	return 0, nil
}

/*
Write implements the io.Writer interface.
*/
func (e *Entity) Write(p []byte) (n int, err error) {
	return 0, nil
}

/*
Close implements the io.Closer interface.
*/
func (e *Entity) Close() error {
	return nil
}
