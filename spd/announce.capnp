using Go = import "/go.capnp";

@0xf454c62f08bc504b;

$Go.package("spd");
$Go.import("spd");

interface Announce {
    public @0 (datagram :Data) -> (crowd :List(Data));
}

