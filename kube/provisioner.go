package kube

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/brazil"
	"github.com/theapemachine/wrkspc/errnie"
	"gopkg.in/yaml.v2"
)

/*
Provisioner brings clusters up.
*/
type Provisioner struct {
	cfg            Config
	client         RestClient
	extendedClient ExtendedClient
	ebsclient      OpenEBSClient
	config         string
	rootpath       string
	stack          []MigratableKind
}

/*
Kind is a wrapper around the file data of a manifest.
*/
type Kind struct {
	Kind string `yaml:"kind"`
	Path string
	File []byte
	Line string
}

/*
NewProvisioner constructs an object that can be used to configure and run
full deployment plan, designed to bring up a Kubernetes cluster.
*/
func NewProvisioner(config string) Provisioner {
	errnie.Traces()
	cfg := NewConfig()

	return Provisioner{
		cfg:            cfg,
		client:         NewRestClient(),
		extendedClient: NewExtendedClient(cfg),
		ebsclient:      NewOpenEBSClient(cfg),
		config:         config,
		stack:          make([]MigratableKind, 0),
	}
}

/*
Deploy reads in the cluster configuration referenced by name in the ~/.wrkspc.yml
file and deploys an exact version of that configuration. It keeps an internal
stack of deployed cluster Kinds so it can be gracefully torn down.
*/
func (infra Provisioner) Deploy() Provisioner {
	errnie.Traces()
	plan := viper.GetStringSlice("wrkspc.kube.distro")
	infra.rootpath = brazil.Workdir() + "/manifests/kubernetes/"

	// Let's bring the entire cluster up.
	for _, step := range plan {
		// Append the new cluster Kind to the provisioner stack so we have a
		// reference should we want to tear things back down if we have too
		// many critical errors during up.
		for _, fi := range brazil.ReadPath(infra.rootpath + step) {
			if fi.IsDir() {
				continue
			}
			path := infra.rootpath + step + "/" + fi.Name()
			errnie.Logs("open file ", path).With(errnie.INFO)
			infra.stack = append(infra.stack, infra.nextStep(path, false)...)
		}
	}

	return infra
}

/*
Teardown takes the entire cluster down piece by piece in reverse order than it
was built. It is theoretically the best way to gracefully destroy a cluster.
TODO: Save the provisioner stack to a yml file so it can be used after the fact.
      Don't be stupid, just load the initial config and reverse it, duh.
*/
func (infra Provisioner) Teardown() Provisioner {
	errnie.Traces()
	plan := viper.GetStringMapStringSlice("infrastructures." + infra.config)

	// Load the plan in the stack format so we can reverse it.
	for _, step := range plan["config"] {
		// Append the new cluster Kind to the provisioner stack so we have a
		// reference should we want to tear things back down if we have too
		// many critical errors during up.

		for _, fi := range brazil.ReadPath(infra.rootpath + step) {
			infra.stack = append(
				infra.stack,
				infra.nextStep(infra.rootpath+step+"/"+fi.Name(), false)...,
			)
		}
	}

	return infra
}

func (infra Provisioner) nextStep(step string, direction bool) []MigratableKind {
	errnie.Traces()
	switch filepath.Ext(step) {
	case ".yml", ".yaml":
		// Since many manifest files are built up of individually creatable items
		// we need to separate those and collect them in a sub stack that we can
		// merge with the overal stack that describes the deployment plan upstream.
		subs := infra.splitTypes(step, "---")
		subStack := make([]MigratableKind, len(subs))

		for idx, t := range subs {
			kind, _ := infra.getKind(step, t, true)

			// Guard against manifests that have a split separator at the end
			// of the file without any content below it. This is valid but
			// requires an extra check.
			if len(kind.File) == 0 {
				continue
			}

			// Add a single creatable item to the sub stack we are building
			// up which is to be merged in the over deployment plan.
			subStack[idx] = infra.provisionForKind(kind, direction)
		}

		return subStack
	case ".sh":
		// Similarly shell scipts will be lines of individual commands we want to
		// pass to our provisioner to execute.
		subs := infra.splitTypes(step, "\n")
		subStack := make([]MigratableKind, len(subs))

		for idx, t := range subs {
			// This needs no unmarshaling, so pass false as the final argument.
			kind, _ := infra.getKind(step, t, false)

			// Guard against empty lines.
			if kind.Line == "" {
				continue
			}

			// Add a single creatable item to the sub stack we are building
			// up which is to be merged in the over deployment plan.
			subStack[idx] = infra.provisionForKind(kind, direction)
		}

		return subStack
	}

	return []MigratableKind{}
}

func (infra Provisioner) provisionForKind(kind Kind, direction bool) MigratableKind {
	errnie.Traces()
	var module MigratableKind

	switch kind.Kind {
	case "Namespace":
		module = NewNamespace(infra.client, kind.File)
	case "ServiceAccount":
		module = NewServiceAccount(infra.client, kind.File)
	case "ClusterRole":
		module = NewClusterRole(infra.client, kind.File)
	case "ClusterRoleBinding":
		module = NewClusterRoleBinding(infra.client, kind.File)
	case "CustomResourceDefinition":
		module = NewCustomResourceDefinition(infra.extendedClient, kind.File)
	case "CSIDriver":
		module = NewCSIDriver(infra.client, kind.File)
	case "PriorityClass":
		module = NewPriorityClass(infra.client, kind.File)
	case "StatefulSet":
		module = NewStatefulSet(infra.client, kind.File)
	case "ConfigMap":
		module = NewConfigMap(infra.client, kind.File)
	case "DaemonSet":
		module = NewDaemonSet(infra.client, kind.File)
	case "Deployment":
		module = NewDeployment(infra.client, kind.File)
	case "Service":
		module = NewService(infra.client, kind.File)
	case "CStorPoolCluster":
		module = NewCStorPoolCluster(infra.ebsclient, kind.File)
	case "PersistentVolumeClaim":
		module = NewPersistentVolumeClaim(infra.client, kind.File)
	case "StorageClass":
		module = NewStorageClass(infra.client, kind.File)
	case "sh":
	// 	module = kube.NewShellExecutor(shell.NewPosh(), kind.Line)
	default:
	}

	if module != nil && direction {
		module.Up()
	}

	return module
}

func (infra Provisioner) splitTypes(step string, seperator string) [][]byte {
	errnie.Traces()
	fd, err := ioutil.ReadFile(filepath.FromSlash(step))
	errnie.Handles(err).With(errnie.KILL)
	return bytes.Split(fd, []byte(seperator))
}

func (infra Provisioner) getKind(step string, fd []byte, isMarshaled bool) (Kind, error) {
	errnie.Traces()
	var kind Kind
	var err error

	if isMarshaled {
		errnie.Logs(string(fd)).With(errnie.DEBUG)
		err := yaml.Unmarshal(fd, &kind)
		errnie.Handles(err).With(errnie.NOOP)
		kind.File = fd
	} else {
		kind.Kind = "sh"
		kind.Line = string(fd)
	}

	kind.Path = step
	return kind, err
}
