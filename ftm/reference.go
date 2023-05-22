package ftm

/*
Reference is a wrapper that mirrors the FollowTheMoney Reference object.
*/
type Reference struct {
	ID         string
	Schema     string
	Properties map[string][]string
}

/*
NewReference creates a new Reference object.
*/
func NewReference() *Reference {
	return &Reference{}
}

/*
Read implements the io.Reader interface.
*/
func (r *Reference) Read(p []byte) (n int, err error) {
	return 0, nil
}

/*
Write implements the io.Writer interface.
*/
func (r *Reference) Write(p []byte) (n int, err error) {
	return 0, nil
}

/*
Close implements the io.Closer interface.
*/
func (r *Reference) Close() error {
	return nil
}
