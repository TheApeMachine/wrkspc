package spd

import (
	"time"

	"github.com/google/uuid"
	"github.com/wrk-grp/errnie"
)

/*
New contructs a Datagram message and marshals it to a byte slice before returning
the object wrapped in a bytes.Buffer.
*/
func New(dataType DataType, role RoleType, scope ScopeType, identity []byte) *Datagram {
	errnie.Trace()

	errnie.Trace()
	var (
		err error
	)

	dg, err := root()

	// Setup the context for the message
	if err != nil {
		errnie.Handles(err)
		return nil
	}

	if err = dg.SetVersion(Version); err != nil {
		errnie.Handles(err)
		return nil
	}

	if err = dg.SetType(dataType); err != nil {
		errnie.Handles(err)
		return nil
	}

	if err = dg.SetRole(role); err != nil {
		errnie.Handles(err)
		return nil
	}

	if err = dg.SetScope(scope); err != nil {
		errnie.Handles(err)
		return nil
	}

	if err = dg.SetIdentity([]byte(identity)); err != nil {
		errnie.Handles(err)
		return nil
	}

	if err = dg.SetUuid([]byte(uuid.New().String())); err != nil {
		errnie.Handles(err)
		return nil
	}

	dg.SetTimestamp(time.Now().UnixNano())
	dg.SetPtr(0)

	return dg
}
