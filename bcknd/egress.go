package bcknd

/*
Egress serves the responses back that were requested through Ingress.
*/
type Egress struct {
}

/*
NewEgress returns a reference to an instance of Egress.
*/
func NewEgress() *Egress {
	return &Egress{}
}
