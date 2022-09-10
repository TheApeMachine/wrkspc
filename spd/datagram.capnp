using Go = import "/go.capnp";
@0x85d3acc39d94e0f8;
$Go.package("spd");
$Go.import("spd/datagram");

struct Datagram {
    uuid      @6 :Text;
    version   @0 :Text;
    role      @1 :Text;
    scope     @2 :Text;
    identity  @3 :Text;
    timestamp @4 :Int64;
    layers    @5 :List(Data);
}
