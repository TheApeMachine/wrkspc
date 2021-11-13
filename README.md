# wrkspc

> This is very much a work in progress.
> Exactly zero of this is to be considered "stable" or has any intention to be so for the moment

Workspace tries to be a single binary "full-stack". It combines infrastructure, clustering,
containerization, services, CI/CD pipelines, and development environment into a singular
end-to-end experience.

Fork it, make it your own. It has an `Unlicense`.

 Features

[x] Built-in Kubernetes Distribution:
    It can deploy itself as a cluster or connect to an existing cluster so all your tools,
	services, and anything else can run on Kubernetes.
 [x] Deploy local Kind cluster when no existing `~/.kube/config` is found.
 [x] Connect to existing cluster is config file is found.
 [x] Include selection of open source package manifests:
		- [x] Weavenet & Weave Scope
		- [x] MetalLB
		- [x] OpenEBS & Min IO
		- [x] Istio & Kiali
		- [x] ArgoCD & Argo Rollouts
		- [x] Argo Workflows & Argo Events
		- [x] Prometheus & Grafana
		- [x] ECK operator & ElasticSearch & Kibana
	- [ ] Deploy production grade cluster.
- [ ] Automatic Replication:
      It can scan any network it is connected to for live hosts in a variety of ways, and use
			the most common ways to try to get access. If it does, it will provision that machine to
			join the cluster.
	- [x] Simple connection based IP scanner.
	- [x] Establish connection over SSH with keyfile or username/password.
	- [ ] Use built-in kubeadm to remotely provision node.
- [x] Built-in Container Runtime:
      There should be no need to pre-install anything. Containerd is integrated as a package and
			a daemon will be started when needed.
	- [ ] Commands are proxied through binary to run in containers.
	- [ ] Commands are automatically aliased on subsequent runs.
- [x] Dynamically Building Development Environment:
      All commands on the terminal are proxied through `wrkspc`, whether they exist or not on
	the local machine. A containerized version of that command (if you first create that image
	of course) is pulled from a registry and executed that way.
[x] stck:
    Stack is a connection with the Docker and Kubernetes APIs through which we can build and
	run Dockerfiles, and deploy yaml based manifest files respectively.
	manifest files. The local path `~/.wrkspc/stck` is created automatically and the default
	config of the cluster distribution is written in there. Anything added or removed from
			that path will be reflected in the cluster.
- [ ] bcknd:
      Backend is a fully featured ETL pipeline that is highly opinionated. It uses an approach
			based on a data lake, schemaless data, no databases, a mono type request path, streams,
			dynamic projections, and a single endpoint. It should cover a lot of use-cases.
	- [ ]


## How it Works

In the examples below you can always run any first command such as `wrkspc` from the path that
contains the source code as `go run main.go`, those two things are equal.

The main idea behind `wrkspc`, besides having your environment dynamically sourced, is to make
code and infrastructure the same. Not `Infrastructure as Code`, but `infra == code`.

Let's look at the functionality from the development side upwards.

## Step 1: Preparing a New Machine

Install (or build) the binary for the platform of the machine you are working on.

For convenience make sure it is somewhere in your `$PATH` so you do not need to use `./`.

Run `wrkspc --version` to get the version and also write the default config file to your home
path. You can have a look at `~/.wrkspc.yml` to see what you can and need to configure.

Then if you don't need anything else than the premade tools in the `./dockerfiles` path of this
repository, there is nothing left to do.

It is more likely though that you will need to wrap whatever tools you use in a Dockerfile and
push them to a registry.

These do not have to be tools that actually exist, for instance the included `pxedust` is a
virtual command that just links to the Dockerfile which uses many elements to build up a full
PXE boot provider.

Done.

I call it "the last fresh install you'll ever need to do."

## Step 2: Development

Depending on how you want to approach things, you can either clear out, or customize `bcknd`
and work in there to have your code integrate more deeply with `wrkspc` and its features.

The benefit of the above method is that everything ultimately remains "inside" the single
binary distribution.

Another approach would be to see your custom services similar to the open source packages
included and/or added later. You can develop and build your services in the usual way and
include their deployment manifest.

Included with this distribution is an opinionated `zsh` configuration and similar one for `vim`.

To get started quickly you can just run the commands below.

```zsh
# Open a new go file for editing.
$ wrkspc run vim main.go
```

It will pull the `zsh` and `vim` container images if they are not already present locally, build,
and then run them.

A tty will be attached when needed.

> An alias is automatically written to your dotfiles and sourced for any command upon their first
> use. That means that from that moment you can omit the `wrkspc run` part for that command.

You can either customize the included configs, or build your own entirely.

## Step 3: Running Services

Using the features of the Cobra CLI package `wrkspc` can run in various "modes".

We have already seen a glimpse of the tools on the terminal in the development mode.

Using the command below we can run a service which will be deployed to a cluster in the
`$USER-dev` namespace.

Let's try this with the built-in backend.

```zsh
$ wrkspc serve bcknd
```

You should now be able to store data.

```zsh
curl --location --request POST 'http://ingress.cluster.local' \
--header 'Content-Type: application/json' \
--data-raw '{
    "context": {
        "role": "datapoint",
				"type": "json",
        "annotations": [{
            "store": "test/myfirstartifact"
        }]
    },
    "data": {
        "payload": []BYTESBUFFER
    }
}'
```

The body of this request carries a specific structured format called a (cloud) datagram, which
can be seen as "one unit" of data traveling through the `bcknd` pipelines.

The `context` header holds meta data that describes the binary encoded `payload`.

The full `datagram` itself can also encode itself to a bytes buffer so it can be stored in most
storage solutions, which in the default setup is `Min IO` or `AWS S3`.

That means that you can always predict the first type when retrieving stored bytes from your
storage solution, which when unmarshaled back to a `datagram` will reveal the inner type through
the `context` header which you can use to unmarshal the `payload`.

 Step 4: Deployment

 you are using all or mostly defaults, some form of cluster is already running. Either it was
ilt as a Kind (Kubernetes in Docker) cluster locally, or an existing config was found and used
 run the development environment.

ere are a few ways in which you can deploy a more production ready version.

```zsh
# Automatic Network Replication will scan any network you are connected to for live hosts and
# open port 22, to which it tries to connect based on the values in `~/.wrkspc.yml`. It will
# provision any connected machine as part of the cluster (as is, nothing is wiped).
$ wrkspc deploy --auto
```

Behind the scenes many things will be moving around to hopefully come to a result of a high
grade cluster setup, with all deployments and services present and working, while also replicated
into (by default) `staging`, `qa`, `production`, and `hotfix` environments.

> A brief mention about "modes". Since wrkspc is not only infrastructure and development environment,
> but also its own services, the built-in cli has commands available to run in the context of these
> modes. While the `deploy --auto` command abstracts that away in automation, you can have a look at
> the `serve` command documentations in `./docs/` to get an understanding on how to deploy a built-in
> service manually.

## Step 5: Operations

By default you get a monitoring and logging stack, for which below is a list of frontends you
can visit.

- [Grafana](https://grafana.cluster.local)
- [Kibana](https://kibana.cluster.local)
- [Kiali](https://kiali.cluster.local)
- [Weave Scope](https://weave.cluster.local)

> Fun fact, there is a couple of custom Grafana plugins included in this repository, so they can
> act as templates, should you be looking to make some.

You can find everything you need should you want to customize things to your needs in the
`~/.wrkspc/` path and the `~/wrkspc.yml`.

## Step 6: Cleanup

If you want to leave the machine you are working on in a clean state it was previously in, you
can use the command below.

```zsh
$ wrkspc cleanup
```

It will delete the whole `~/.wrkspc` path and the `~/.wrkspc.yml` config, as well as the binary
itself, undoing all other changes that were made to the system.

There is a milder version of this command listed below which just removes anything that is not
being used at the time and does some other small tasks to just bring your workspace back into
shape as much as possible.

```zsh
$ wrkspc tidy
```

## Frequently Asked Questions

- So it is monolithic?
  It is yes. But only from a code perspective I suppose. I does break the "no shared resources" idea
	that I have heard about in the context of micro-services, but that was always more an issue still
	under discussion than a solidified solution. The question is do we duplicate code and functionality
	or run the risk of a change in one service breaking another.

	This project runs in "modes" and can therefor act as any number of micro-services in deployment.
	It retains the benefit of one service going down not breaking any other services currently running.

	I see it more as a project with sub-modules without the dirty git sub-module workflow.
	Everything is packages though, so if you needed to, it would be trivial to split them into
	separate repositories.
