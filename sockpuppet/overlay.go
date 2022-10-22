package sockpuppet

/*
Overlay is a wrpapper around Nebula, based on WireGuard.
It provides a zero-trust overlay network to connect any
machine to any other machine, even across subnets.
*/
type Overlay struct{}

/*
NewOverlay constructs a Nebula overlay network based on
WireGuard, which allows a hybrid-cloud setup, connecting
multiple cloud providers and on-prem setups together in
a single cluster configuration.
*/
func NewOverlay() *Overlay {
	return &Overlay{}
}
