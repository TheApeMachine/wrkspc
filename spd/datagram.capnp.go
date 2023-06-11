// Code generated by capnpc-go. DO NOT EDIT.

package spd

import (
	capnp "capnproto.org/go/capnp/v3"
	text "capnproto.org/go/capnp/v3/encoding/text"
	schemas "capnproto.org/go/capnp/v3/schemas"
)

type Datagram capnp.Struct

// Datagram_TypeID is the unique identifier for the type Datagram.
const Datagram_TypeID = 0x9b4169eff5c9b6f0

func NewDatagram(s *capnp.Segment) (Datagram, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 8})
	return Datagram(st), err
}

func NewRootDatagram(s *capnp.Segment) (Datagram, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 8})
	return Datagram(st), err
}

func ReadRootDatagram(msg *capnp.Message) (Datagram, error) {
	root, err := msg.Root()
	return Datagram(root.Struct()), err
}

func (s Datagram) String() string {
	str, _ := text.Marshal(0x9b4169eff5c9b6f0, capnp.Struct(s))
	return str
}

func (s Datagram) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (Datagram) DecodeFromPtr(p capnp.Ptr) Datagram {
	return Datagram(capnp.Struct{}.DecodeFromPtr(p))
}

func (s Datagram) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s Datagram) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s Datagram) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s Datagram) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s Datagram) Checksum() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return []byte(p.Data()), err
}

func (s Datagram) HasChecksum() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s Datagram) SetChecksum(v []byte) error {
	return capnp.Struct(s).SetData(0, v)
}

func (s Datagram) Uuid() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return []byte(p.Data()), err
}

func (s Datagram) HasUuid() bool {
	return capnp.Struct(s).HasPtr(1)
}

func (s Datagram) SetUuid(v []byte) error {
	return capnp.Struct(s).SetData(1, v)
}

func (s Datagram) Version() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(2)
	return []byte(p.Data()), err
}

func (s Datagram) HasVersion() bool {
	return capnp.Struct(s).HasPtr(2)
}

func (s Datagram) SetVersion(v []byte) error {
	return capnp.Struct(s).SetData(2, v)
}

func (s Datagram) Type() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(3)
	return []byte(p.Data()), err
}

func (s Datagram) HasType() bool {
	return capnp.Struct(s).HasPtr(3)
}

func (s Datagram) SetType(v []byte) error {
	return capnp.Struct(s).SetData(3, v)
}

func (s Datagram) Role() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(4)
	return []byte(p.Data()), err
}

func (s Datagram) HasRole() bool {
	return capnp.Struct(s).HasPtr(4)
}

func (s Datagram) SetRole(v []byte) error {
	return capnp.Struct(s).SetData(4, v)
}

func (s Datagram) Scope() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(5)
	return []byte(p.Data()), err
}

func (s Datagram) HasScope() bool {
	return capnp.Struct(s).HasPtr(5)
}

func (s Datagram) SetScope(v []byte) error {
	return capnp.Struct(s).SetData(5, v)
}

func (s Datagram) Identity() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(6)
	return []byte(p.Data()), err
}

func (s Datagram) HasIdentity() bool {
	return capnp.Struct(s).HasPtr(6)
}

func (s Datagram) SetIdentity(v []byte) error {
	return capnp.Struct(s).SetData(6, v)
}

func (s Datagram) Timestamp() int64 {
	return int64(capnp.Struct(s).Uint64(0))
}

func (s Datagram) SetTimestamp(v int64) {
	capnp.Struct(s).SetUint64(0, uint64(v))
}

func (s Datagram) Ptr() int32 {
	return int32(capnp.Struct(s).Uint32(8))
}

func (s Datagram) SetPtr(v int32) {
	capnp.Struct(s).SetUint32(8, uint32(v))
}

func (s Datagram) Layers() (capnp.DataList, error) {
	p, err := capnp.Struct(s).Ptr(7)
	return capnp.DataList(p.List()), err
}

func (s Datagram) HasLayers() bool {
	return capnp.Struct(s).HasPtr(7)
}

func (s Datagram) SetLayers(v capnp.DataList) error {
	return capnp.Struct(s).SetPtr(7, v.ToPtr())
}

// NewLayers sets the layers field to a newly
// allocated capnp.DataList, preferring placement in s's segment.
func (s Datagram) NewLayers(n int32) (capnp.DataList, error) {
	l, err := capnp.NewDataList(capnp.Struct(s).Segment(), n)
	if err != nil {
		return capnp.DataList{}, err
	}
	err = capnp.Struct(s).SetPtr(7, l.ToPtr())
	return l, err
}

// Datagram_List is a list of Datagram.
type Datagram_List = capnp.StructList[Datagram]

// NewDatagram creates a new list of Datagram.
func NewDatagram_List(s *capnp.Segment, sz int32) (Datagram_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 8}, sz)
	return capnp.StructList[Datagram](l), err
}

// Datagram_Future is a wrapper for a Datagram promised by a client call.
type Datagram_Future struct{ *capnp.Future }

func (p Datagram_Future) Struct() (Datagram, error) {
	s, err := p.Future.Struct()
	return Datagram(s), err
}

const schema_85d3acc39d94e0f8 = "x\xdaL\xce\xbd\x8a\xd4P\x1c\x05\xf0snn&\x19" +
	"p\xd6\x0d\xf7\x0a+\xac\xf8\xc1\x0a\x8b\xb8\xba+\x8a\xba" +
	"\x08~`e5\xe3}\x82K\x12\xd6\xe0f&$w" +
	"\x84\xa9\xac|\x03;\xf1\x0d\xacm-\xc4\xc2\xc2B|" +
	"\x81\xc5\x17\xd0F\xb0\x8b\xfc\x03\x8e\x96\xe7w\xce\xbd\xfc" +
	"7O\x1e\xe8\x83\xc9\x07B\xcdl<\xea\x7f\xbe\xff\xfc" +
	"\xebG\xf5\xf0\x0df\x86\xaa\xff}\xf2\xfa\xed\xc7w\xdf" +
	"^!N\x13 \xfb\xf4=\xfb\x9a\x00\x07_n+\xec" +
	"\xf5]S\\/|\xf0\xea\xa8\xf5\xf5\xb5\xdc7\xf3\xe6" +
	"\xf0\xb1\x0f\xfe(i}=%g\xbb\x91\x064\x013" +
	"\xe6\x13\xc0\xa5\x8c\xe8,\x153\xd2R<\xe3\x15\xc0\x9d" +
	"\x12\xdf\x12W\xcaR\x01\xe6\x0c\x1f\x01nS|[<" +
	"\x8a,#\xc0\x9c\x1d\xf6V\xfc\x82\xb8\xd6\x96\x1a0\xe7" +
	"\x06\xdf\x12\xdf\x11\x8fc\xcb\x180\x17y\x03p\xdb\xe2" +
	"\xbb\xe2\xa3\x91\xe5\x080\x97\x87{v\xc4\xf7\xa9\xc8\xc4" +
	"2\x01\xcc\x1e\x9f\x02\xee\xaa\xf0\x1d\x99\xa7\xca2\x05\xcc" +
	"-^\x02\xdc\xbe\xf8=\xf1qb9\x06\xcc]\x1e\x02" +
	"\xee\xa6\xf8\x94\x8a}\xfe\xac\xcc\x9fw\xcb\x1a\x00'P" +
	"\x9c\x80\xa7\x97\xcb\xaa\xf8\x1b^\xbe(\xdb\xaeZ\xcc\xd7" +
	"eX5\xe5:\xb4\x8b\xe3u8\xdf\xe5\x8b\x7fU_" +
	"\x15\xe5<Ta\xf5\xdf\xc7}\xa8\xea\xb2\x0b\xbe\x06\x1b" +
	"\xc6P\x8c\xc1\xa4\x09-5\x145x\xff\xd8\xaf\xca\xb6" +
	"\xe3\x068\x8d8\xbc\xda\x00\xff\x04\x00\x00\xff\xff\xaf\xe4" +
	"N\xa4"

func init() {
	schemas.Register(schema_85d3acc39d94e0f8,
		0x9b4169eff5c9b6f0)
}
