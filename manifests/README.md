Manifests

 this directory you will find all the Kubernetes manifests, Dockerfiles, and source code included as part of the distribution of `wrkspc`.

brief description of the main items as sub headers below.

## Dockerfiles

ntains opinionated tool configurations wrapped as Dockerfiles so they can be used by the workspace command proxy.

 Kubernetes

ntains pre configured yaml files for various open source services that build up the basic stack for this distribution.

### TODO
- [ ] Patch Strict ARP for Weavenet.
- [ ] Patch OpenEBS default storage class.
- [ ] Download weave binary

## Bcknd

Contains the Go source code for the built in backend ETL pipeline.
