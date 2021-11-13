package contempt

/*
Range is a simple type that describes a linear range between two values. It is used mostly to keep
input types small in count.
*/
type Range struct {
	From    int
	To      int
	current int
}

/*
Next updates the `current` field and returns its value.
*/
func (r *Range) Next() int {
	r.current = r.From + r.current
	return r.current
}
