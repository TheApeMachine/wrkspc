using Go = import "/go.capnp";
@0x85d3acc39d94e0f8;
$Go.package("spd");
$Go.import("spd/datagram");

struct Datagram {
    uuid      @6 :Data;
    version   @0 :Data;
    type      @7 :Data;
    role      @1 :Data;
    scope     @2 :Data;
    identity  @3 :Data;
    timestamp @4 :Int64;
    ptr       @8 :Int32;
    layers    @5 :List(Data);
}
