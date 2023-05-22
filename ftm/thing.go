package ftm

import "io"

/*
Thing defines a common interface for all FollowTheMoney objects.
*/
type Thing interface {
	io.ReadWriteCloser
}
