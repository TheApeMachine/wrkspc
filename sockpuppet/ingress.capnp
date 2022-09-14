using Go = import "/go.capnp";

@0xf454c62f08bc504b;

$Go.package("sockpuppet");
$Go.import("sockpuppet");

interface Ingress {
	handler @0 (datagram :Data) -> (data :Data);
}
