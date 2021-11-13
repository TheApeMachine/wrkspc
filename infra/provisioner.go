package infrastructure

import (
	"bytes"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/theapemachine/wrkspc/errnie"
	"github.com/theapemachine/wrkspc/kube"
	"gopkg.in/yaml.v2"
)

type Provisioner struct {
	cfg            kube.Config
	client         kube.RestClient
	extendedClient kube.ExtendedClient
	ebsclient      kube.OpenEBSClient
	config         string
	rootpath       string
	stack          []kube.MigratableKind
}

/*
KubeKind is a wrapper around the file data of a manifest.
*/
type KubeKind struct {
	Kind string `yaml:"kind"`
	Path string
	File []byte
	Line string
}

/*
NewwProvisioner constructs an object that can be used to configure and run
full deployment plan, designed to bring up a Kubernetes cluster.
*/
func NewProvisioner(config string) Provisioner {
	cfg := kube.NewConfig()

	return Provisioner{
		cfg:            cfg,
		client:         kube.NewRestClient(),
		extendedClient: kube.NewExtendedClient(cfg),
		ebsclient:      kube.NewOpenEBSClient(cfg),
		config:         config,
		stack:          make([]kube.MigratableKind, 0),
	}
}

/*
Deploy reads in the cluster configuration referenced by name in the ~/.wrkspc.yml
file and deploys an exact version of that configuration. It keeps an internal
stack of deployed cluster Kinds so it can be gracefully torn down.
*/
func (infra Provisioner) Deploy() error {
	plan := viper.GetStringMapStringSlice("infrastructures." + infra.config)
	infra.rootpath = plan["rootpath"][0]

	// Let's bring the entire cluster up.
	for _, step := range plan["config"] {
		// Append the new cluster Kind to the provisioner stack so we have a
		// reference should we want to tear things back down if we have too
		// many critical errors during up.
		infra.stack = append(infra.stack, infra.nextStep(infra.rootpath+step, true)...)
	}

	return nil
}

/*
Teardown takes the entire cluster down piece by piece in reverse order than it
was built. It is theoretically the best way to gracefully destroy a cluster.
TODO: Save the provisioner stack to a yml file so it can be used after the fact.
      Don't be stupid, just load the initial config and reverse it, duh.
*/
func (infra Provisioner) Teardown() error {
	plan := viper.GetStringMapStringSlice("infrastructures." + infra.config)

	// Load the plan in the stack format so we can reverse it.
	for _, step := range plan["config"] {
		// Append the new cluster Kind to the provisioner stack so we have a
		// reference should we want to tear things back down if we have too
		// many critical errors during up.
		infra.stack = append(infra.stack, infra.nextStep(infra.rootpath+step, false)...)
	}

	// Tear it all back down in reverse order.
	for i := len(infra.stack) - 1; i > 0; i-- {
		// Guard against modules that did not come up and produced a nil value
		// in the provisioner reference stack.
		if infra.stack[i] == nil {
			continue
		}

		// Call the down method in the MigratableKind.
		if err := infra.stack[i].Down(); err != nil {
			errnie.Handles(err).With(errnie.KILL)
		}
	}

	return nil
}

func (infra Provisioner) nextStep(step string, direction bool) []kube.MigratableKind {
	switch filepath.Ext(step) {
	case ".yml", ".yaml":
		// Since many manifest files are built up of individually creatable items
		// we need to separate those and collect them in a sub stack that we can
		// merge with the overal stack that describes the deployment plan upstream.
		subs := infra.splitTypes(step, "---")
		subStack := make([]kube.MigratableKind, len(subs))

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
		subStack := make([]kube.MigratableKind, len(subs))

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

	return []kube.MigratableKind{}
}

func (infra Provisioner) provisionForKind(kind KubeKind, direction bool) kube.MigratableKind {
	var module kube.MigratableKind

	switch kind.Kind {
	case "Namespace":
		module = kube.NewNamespace(infra.client, kind.File)
	case "ServiceAccount":
		module = kube.NewServiceAccount(infra.client, kind.File)
	case "ClusterRole":
		module = kube.NewClusterRole(infra.client, kind.File)
	case "ClusterRoleBinding":
		module = kube.NewClusterRoleBinding(infra.client, kind.File)
	case "CustomResourceDefinition":
		module = kube.NewCustomResourceDefinition(infra.extendedClient, kind.File)
	case "CSIDriver":
		module = kube.NewCSIDriver(infra.client, kind.File)
	case "PriorityClass":
		module = kube.NewPriorityClass(infra.client, kind.File)
	case "StatefulSet":
		module = kube.NewStatefulSet(infra.client, kind.File)
	case "ConfigMap":
		module = kube.NewConfigMap(infra.client, kind.File)
	case "DaemonSet":
		module = kube.NewDaemonSet(infra.client, kind.File)
	case "Deployment":
		module = kube.NewDeployment(infra.client, kind.File)
	case "Service":
		module = kube.NewService(infra.client, kind.File)
	case "CStorPoolCluster":
		module = kube.NewCStorPoolCluster(infra.ebsclient, kind.File)
	case "PersistentVolumeClaim":
		module = kube.NewPersistentVolumeClaim(infra.client, kind.File)
	case "StorageClass":
		module = kube.NewStorageClass(infra.client, kind.File)
	case "sh":
	// 	module = kube.NewShellExecutor(shell.NewPosh(), kind.Line)
	default:
	}

	if module != nil && direction {
		errnie.Handles(module.Up()).With(errnie.KILL)
	}

	return module
}

func (infra Provisioner) splitTypes(step string, seperator string) [][]byte {
	fd, err := ioutil.ReadFile(filepath.FromSlash(step))
	errnie.Handles(err).With(errnie.KILL)
	return bytes.Split(fd, []byte(seperator))
}

func (infra Provisioner) getKind(step string, fd []byte, isMarshaled bool) (KubeKind, error) {
	var kind KubeKind
	var err error

	if isMarshaled {
		err := yaml.Unmarshal(fd, &kind)
		errnie.Handles(err).With(errnie.KILL)
		kind.File = fd
	} else {
		kind.Kind = "sh"
		kind.Line = string(fd)
	}

	kind.Path = step
	return kind, err
}
