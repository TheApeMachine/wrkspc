package matroesjka

import (
	"embed"
	"syscall"
	"unsafe"

	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
Embed a mini filesystem into the binary with the contents of ./bin so we can use it to run any
external dependencies right from an in-memory filesystem.
*/
//go:embed bin/*
var payload embed.FS

/*
Embed is a wrapper around embedded executable binaries.
*/
type Embed struct {
	name string
}

/*
NewEmbed prepares an embedded binary for in-memory execution.
*/
func NewEmbed(name string) *Embed {
	errnie.Traces()
	return &Embed{name: name}
}

/*
Write out the embedded dependencies.
*/
func (embedfs *Embed) Write() {
	brazil.MakePath(brazil.HomePath() + "/wrkspc")
	brazil.WriteIfNotExists("bin/runc", payload, true)
	brazil.WriteIfNotExists("bin/containerd", payload, true)
	brazil.WriteIfNotExists("bin/containerd-shim-runc-v2", payload, true)
	// brazil.WriteIfNotExists("bin/dockerd", payload, true)
	// brazil.WriteIfNotExists("bin/docker-proxy", payload, true)
	// brazil.WriteIfNotExists("bin/docker-init", payload, true)
	// brazil.WriteIfNotExists("bin/docker", payload, true)
	brazil.WriteIfNotExists("bin/modprobe", payload, true)
}

/*
Exec ...
*/
func (embedfs *Embed) Exec() {
	errnie.Traces()

	fd, err := embedfs.MemfdCreate("/runc.bin")
	errnie.Handles(err).With(errnie.KILL)
	errnie.Logs(fd).With(errnie.DEBUG)
	errnie.Handles(
		embedfs.CopyToMem(
			fd, []byte{},
		),
	).With(errnie.KILL)
	errnie.Handles(embedfs.ExecveAt(fd)).With(errnie.KILL)
}

/*
MemfdCreate creates an in-memory file descriptor.
*/
func (embedfs *Embed) MemfdCreate(path string) (r1 uintptr, err error) {
	errnie.Traces()

	s, err := syscall.BytePtrFromString(path)
	if !errnie.Handles(err).With(errnie.KILL).OK {
		return 0, err
	}

	errnie.Logs(s).With(errnie.DEBUG)

	r1, _, errno := syscall.Syscall(319, uintptr(unsafe.Pointer(s)), 0, 0)

	errnie.Logs(r1).With(errnie.DEBUG)

	if int(r1) == -1 {
		return r1, errno
	}

	return r1, nil
}

/*
CopyToMem writes the embedfs file into a memory location.
*/
func (embedfs *Embed) CopyToMem(fd uintptr, buf []byte) (err error) {
	errnie.Traces()

	_, err = syscall.Write(int(fd), buf)
	errnie.Handles(err).With(errnie.KILL)
	return err
}

/*
ExecveAt executes the in-memory binary.
*/
func (embedfs *Embed) ExecveAt(fd uintptr) (err error) {
	errnie.Traces()

	s, err := syscall.BytePtrFromString("")
	errnie.Handles(err).With(errnie.KILL)
	errnie.Logs(s).With(errnie.DEBUG)

	ret, _, errno := syscall.Syscall6(322, fd, uintptr(unsafe.Pointer(s)), 0, 0, 0x1000, 0)
	if int(ret) == -1 {
		return errno
	}

	// never hit
	errnie.Logs("you should not be here").With(errnie.ERROR)
	return err
}
