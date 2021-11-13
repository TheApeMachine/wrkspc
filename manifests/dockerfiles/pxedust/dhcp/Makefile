make setup:
	wget https://stable.release.flatcar-linux.net/amd64-usr/current/flatcar_production_pxe.vmlinuz
	wget https://stable.release.flatcar-linux.net/amd64-usr/current/flatcar_production_pxe.vmlinuz.sig
	wget https://stable.release.flatcar-linux.net/amd64-usr/current/flatcar_production_pxe_image.cpio.gz
	wget https://stable.release.flatcar-linux.net/amd64-usr/current/flatcar_production_pxe_image.cpio.gz.sig

build:
	docker build -t theapemachine/pxe:latest .

push:
	docker push theapemachine/pxe:latest

buildpush:
	docker build -t theapemachine/pxe:latest .
	docker push theapemachine/pxe:latest

run:
	docker run --rm --net=host --cap-add=NET_ADMIN theapemachine/pxe:latest

shell:
	docker run --rm -ti --net=host --cap-add=NET_ADMIN theapemachine/pxe:latest /bin/bash

validate:
	cat drknow-ignite-boot.yml | docker run --rm -i quay.io/coreos/ct:latest-dev --platform=linux
	cat drknow-ignite.yml | docker run --rm -i quay.io/coreos/ct:latest-dev --platform=linux
