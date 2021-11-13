package bcknd

/*
Ingress is the only "endpoint" that this backend architecture has.
I accepts only the Datagram format, but wrapped in that type you
can store any other kind of data.
*/
type Ingress struct {
}

/*
NewIngress returns a reference to an instance of Ingress.
*/
func NewIngress() *Ingress {
	return &Ingress{}
}
