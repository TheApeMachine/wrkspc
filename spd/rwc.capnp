using Go = import "/go.capnp";

@0xf454c62f08bc504b;

$Go.package("spd");
$Go.import("spd");

interface RWC {
    read  @0 () -> (out :Data);
    write @1 (in :Data);
}

