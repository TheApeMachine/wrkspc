package matrix

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
)

/*
RootFS is a wrapper around a .tar.gz root filesystem defined as a base
image in the current config file.
*/
type RootFS struct {
	root  string
	tool  string
	osmap map[string][]string
}

/*
NewRootFS gets a handle on a special builder that writes a root file system as tar file to the
root of a new (scratch) container image.
*/
func NewRootFS(root string, tool string) RootFS {
	errnie.Traces()

	return RootFS{
		root:  root,
		tool:  tool,
		osmap: viper.GetStringMapStringSlice("wrkspc.atomic"),
	}
}

/*
Build the root filesystem and write it to the container image.
*/
func (rootfs RootFS) Build() {
	errnie.Traces()

	// Match the base os to the tool we are building.
	rootos := rootfs.lookupOS()

	if viper.GetBool("debug") {
		wd, err := os.Getwd()
		errnie.Handles(err).With(errnie.KILL)
		rootfs.root = filepath.FromSlash(wd + "/manifests/images")
	}

	// Copy the root filesystem of the base os to the build context.
	brazil.Copy(
		filepath.FromSlash(rootfs.root+"/"+rootos+".tar.gz"),
		filepath.FromSlash(rootfs.root+"/rootfs.tar.gz"),
	)

	// Since this Build method is called on any other build and is
	// therefor recursive, we need to specify the correct tool to
	// build, as well as tell the build flow to skip some steps.
	buildflow := NewBuild("rootfs")
	buildflow.Atomic(true) // Make the recursive call to Atomic and pass true this time to skip the
	// rebuilding of the root file system.
}

func (rootfs RootFS) lookupOS() string {
	errnie.Traces()

	for oskey, tools := range rootfs.osmap {
		if rootos := rootfs.iterTools(oskey, tools); rootos != "" {
			return rootos
		}
	}

	return ""
}

func (rootfs RootFS) iterTools(oskey string, tools []string) string {
	errnie.Traces()

	for _, tool := range tools {
		if tool == rootfs.tool {
			return oskey
		}
	}

	return ""
}
