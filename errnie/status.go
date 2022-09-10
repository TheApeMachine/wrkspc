package errnie

/*
Status is a simple check to see if the overall system is ok or not.
*/
type Status uint

const (
	// OK indicates nothing is wrong.
	OK Status = iota
	// NOOK indicates something is wrong.
	NOOK
)

/*
Reason is a second layer of granularity over which we can inspect
the state when things are not ok.
*/
type Reason uint

const (
	// ERR indicates an error that came in as not nil.
	ERR Reason = iota
)
