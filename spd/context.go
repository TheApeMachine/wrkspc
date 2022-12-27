package spd

import (
	"time"

	capnp "capnproto.org/go/capnp/v3"
	"github.com/google/uuid"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/tweaker"
)

func root() (*Datagram, error) {
	arena := capnp.SingleSegment(nil)

	_, seg, err := capnp.NewMessage(arena)
	if errnie.Handles(err) != nil {
		return &Empty, err
	}

	dg, err := NewRootDatagram(seg)
	if errnie.Handles(err) != nil {
		return &Empty, err
	}

	return &dg, nil
}

func (dg *Datagram) setContext(
	media MediaType, role RoleType, scope ScopeType,
) error {
	var err error

	if err = dg.SetUuid([]byte(uuid.New().String())); err != nil {
		return errnie.Handles(err)
	}

	if err = dg.SetVersion(Version); err != nil {
		return errnie.Handles(err)
	}

	dg.SetTimestamp(time.Now().UnixNano())

	if err = dg.SetType(media); err != nil {
		return errnie.Handles(err)
	}

	if err = dg.SetRole(role); err != nil {
		return errnie.Handles(err)
	}

	if err = dg.SetScope(scope); err != nil {
		return errnie.Handles(err)
	}

	if err = dg.SetIdentity(tweaker.GetIdentity()); err != nil {
		return errnie.Handles(err)
	}

	return err
}
