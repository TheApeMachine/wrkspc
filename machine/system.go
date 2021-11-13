package machine

/*
System represents a machine and any values we are interested in.
*/
type System struct {
	Ulimit int64
}

/*
NewSystem returns an empty system we can initialize.
*/
func NewSystem() *System {
	system := &System{}
	system.Ulimit = system.getUlimit()
	return system
}

/*
getUlimit gets the maximum allowed open file descriptors of the system.
TODO: This is not really working as expected, hard-coding for now.
*/
func (system *System) getUlimit() int64 {
	// out, err := exec.Command("ulimit", "-n").Output()
	// errnie.Handles(err).With(errnie.KILL)

	// s := strings.TrimSpace(string(out))
	// i, err := strconv.ParseInt(s, 10, 64)
	// errnie.Handles(err).With(errnie.KILL)

	// _ = i

	return 1024
}
