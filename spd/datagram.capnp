using Go = import "/go.capnp";
@0x85d3acc39d94e0f8;
$Go.package("spd");
$Go.import("spd/datagram");

struct Datagram {
	checksum  @0 :Data;
    uuid      @1 :Data;
    version   @2 :Data;
    type      @3 :Data;
    role      @4 :Data;
    scope     @5 :Data;
    identity  @6 :Data;
    timestamp @7 :Int64;
    ptr       @8 :Int32;
    layers    @9 :List(Data);
}