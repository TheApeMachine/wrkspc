package matrix

import (
	"os"
	"path/filepath"
	"strings"

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

func NewRootFS(root string, tool string) RootFS {
	return RootFS{
		root:  root,
		tool:  tool,
		osmap: viper.GetStringMapStringSlice("base-os"),
	}
}

func (rootfs RootFS) Build() {
	// Match the base os to the tool we are building.
	rootos := rootfs.lookupOS()

	if viper.GetBool("debug") {
		wd, err := os.Getwd()
		errnie.Handles(err).With(errnie.KILL)
		rootfs.root = filepath.FromSlash(wd + "/" + "dockerfiles")
	}

	// Copy the root filesystem of the base os to the build context.
	brazil.Copy(
		filepath.FromSlash(rootfs.root+"/images/"+rootos+".tar.gz"),
		filepath.FromSlash(rootfs.root+"/rootfs/rootfs.tar.gz"),
	)

	// Since this Build method is called on any other build and is
	// therefor recursive, we need to specify the correct tool to
	// build, as well as tell the build flow to skip some steps.
	buildflow := NewBuild("rootfs", strings.Split(rootos, "-")[0])
	buildflow.Atomic(true)
}

func (rootfs RootFS) lookupOS() string {
	for oskey, tools := range rootfs.osmap {
		if rootos := rootfs.iterTools(oskey, tools); rootos != "" {
			return rootos
		}
	}

	return ""
}

func (rootfs RootFS) iterTools(oskey string, tools []string) string {
	for _, tool := range tools {
		if tool == rootfs.tool {
			return oskey
		}
	}

	return ""
}
