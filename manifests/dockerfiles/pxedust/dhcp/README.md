# PXE Server

This project provisions a docker container which bundles everything needed to run a PXE boot server.

It is specifically designed to provision machines with [Flatcar Linux](https://www.kinvolk.io/flatcar-container-linux/), a minimal container OS.

Flatcar Linux features:

- Immutable infrastrure
- Designed to scale
- Reduced complexity
- Secure by design
- Security patch automation
- Immutable filesystem
- Minimal attack surface
- Self-driving updates
- Possibility to rollback

Flatcar Linux is a fork of CoreOS, which received an End of Life event in 2020.

It (and the company behind it, Kinvolk) is currently aquired by Microsoft.

The full documentation lives [here](https://kinvolk.io/docs/flatcar-container-linux/latest).

## Things to Know

Any machine that boots, has PXE boot as a selected mechanism, and finds this server, will be entirely wiped and installed with Flatcar Linux.

For now you will have to go through the *-ignite-boot.yml files and manually point the IP to your iPXE host (where you are running the docker container).

Also inside the Dockerfile you will have to make sure to set the IPs to match the nodes you are setting up. This is now partly dynamically configured with build arguments in the Dockerfile.

Currently we follow a bit of a pattern to assign hostnames to machine based on computer systems from fiction (movies, books, etc.) but that is likely to change in the future as it is probably not very scalable.

At such a point that it makes sense we can switch to an automatically generated and assigned hostname system and fully move over to the concept of [treating our systems as cattle and not pets](http://cloudscaling.com/blog/cloud-computing/the-history-of-pets-vs-cattle/).

## Configuration

Flatcar Linux is configured using a yaml format which uses a transpiler mechanism to bring it to its eventual state, json, which is how the OS can be configured to a detail that is familiar in ways you would traditionally install additional functionality, or configure existing features (networking, user accounts, services, etc.)

They will be validated by the transpiler as well to prevent corrupted configuration being applied, and it uses a form of adapter pattern to make them compatible with many end-targets, providers of cloud services for example.

These configs will be included in the container built that serves as the main purpose and output of this repository and can be validated locally for some additional time savings, should something as simple as a typo accidentally go unnoticed.

To do so you can look in the Makefile, as holds true for all following sections, but the common way is shown below.

```bash
$ make validate
```

## How to Build

```bash
$ make build
```

## How to Run

```bash
$ make run
```

Another way to run, which could be used on servers or local machines is to use the `docker-compose.yml`. In that case you will need nothing but that file to bring up a machine.

```bash
$ docker-compose up
```

It will start the PXE server on whatever machine you choose to host it on, then all you need to do is make sure PXE boot is turned on in the bios of machines you want to provision with Flatcar Linux, fully configured by the configuration files that live in `node-configs`.

## How to Push Changes

```bash
$ make push
```

Or a shortcut to build and push.

```bash
$ make buildpush
```
