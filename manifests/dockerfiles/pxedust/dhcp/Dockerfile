FROM ubuntu:18.04

# These values should likely be customized using build arguments during the docker build process.
# They are pretty conventional values though, and if at all possible I would configure the network
# to match these, it will be much simpler to get things running smoothly first try.
ARG MY_PXE_SERVER_IP="192.168.0.164"
ARG MY_DHCP_RANGE="192.168.0.201,192.168.0.230"
ARG MY_GATEWAY_IP="192.168.0.1"
ARG MY_DNS_SERVER_IP="192.168.0.1"

# The tftp root directory we will be serving the netboot images and files from. The vmlinuz and initrd
# compatible images will live here, as well as the grub config.
WORKDIR /srv/tftp

# Install packages we need. Most importantly configuring things such that booting from uefi bios becomes
# possible, since not all network cards have legacy drivers installed.
ARG PACKAGES="dnsmasq=2.79-1 wget=1.19.4-1ubuntu2.2 dpkg-dev=1.19.0.5ubuntu2.3 build-essential git ca-certificates curl"
RUN set -x \
    && sed -i -e 's/^# deb-src /deb-src /' /etc/apt/sources.list \
    && apt-get update -qq \
    && apt-get install -qq --no-install-recommends ${PACKAGES} \
    && wget -qO grubx64.efi http://archive.ubuntu.com/ubuntu/dists/bionic/main/uefi/grub2-amd64/current/grubnetx64.efi.signed \
    && cd /tmp \
    && apt-get source -qq shim-signed \
    && cp shim-signed-*/shimx64.efi.signed /srv/tftp/bootx64.efi \
    && rm -rf /var/lib/apt/lists/*

# We need to setup an http server to serve a connection site that will be called upon during the
# flatcar linux first boot. It will be serving the transpiled configuration files to setup the OS.
RUN git clone https://github.com/emikulic/darkhttpd.git \
    && cd darkhttpd \
    && make && chmod +x darkhttpd && mv darkhttpd /usr/bin

# Copy the grub config which is the main trick to getting uefi boot to work and the dnsmasq configuration
# which provides the PXE boot capabilities to remotely and automatically install the OS images.
COPY grub.cfg /srv/tftp/grub/grub.cfg
COPY dnsmasq.conf /etc/dnsmasq.conf

# Add the netboot OS images to the tftp site.
ADD flatcar_production_pxe_image.cpio.gz /srv/tftp/
ADD flatcar_production_pxe.vmlinuz /srv/tftp/

# Search and replace the template values in the dnsmasq config to use the IPs provided as arguments
# to the build process, reducing some manual overhead.
RUN set -x \
    && sed -i \
    -e "s/MY_DHCP_RANGE/${MY_DHCP_RANGE}/g" \
    -e "s/MY_GATEWAY_IP/${MY_GATEWAY_IP}/g" \
    -e "s/MY_DNS_SERVER_IP/${MY_DNS_SERVER_IP}/g" \
    /etc/dnsmasq.conf \
    && sed -i -e "s/MY_PXE_SERVER_IP/${MY_PXE_SERVER_IP}/g" /srv/tftp/grub/grub.cfg

# Download and setup the config transpiler which does some sanity checks on the yaml 
# configs and converts them to json.
RUN CT_VER=v0.6.1 \
    && ARCH=x86_64 \
    && OS=unknown-linux-gnu \
    && DOWNLOAD_URL=https://github.com/coreos/container-linux-config-transpiler/releases/download \
    && rm -f /tmp/ct-${CT_VER}-${ARCH}-${OS} /tmp/ct-${CT_VER}-${ARCH}-${OS}.asc /tmp/coreos-app-signing-pubkey.gpg \
    && curl -L ${DOWNLOAD_URL}/${CT_VER}/ct-${CT_VER}-${ARCH}-${OS} -o /tmp/ct-${CT_VER}-${ARCH}-${OS} \
    && chmod u+x /tmp/ct-${CT_VER}-${ARCH}-${OS} \
    && curl https://coreos.com/dist/pubkeys/app-signing-pubkey.gpg -o /tmp/coreos-app-signing-pubkey.gpg \
    # && gpg2 --import --keyid-format LONG /tmp/coreos-app-signing-pubkey.gpg \
    && curl -L ${DOWNLOAD_URL}/${CT_VER}/ct-${CT_VER}-${ARCH}-${OS}.asc -o /tmp/ct-${CT_VER}-${ARCH}-${OS}.asc \
    # && gpg2 --verify /tmp/ct-${CT_VER}-${ARCH}-${OS}.asc /tmp/ct-${CT_VER}-${ARCH}-${OS} \
    && mv /tmp/ct-${CT_VER}-${ARCH}-${OS} /usr/bin/ct

# Include the configs to setup and configure the OS and perform the transpile process to convert
# into the final json format needed.
WORKDIR /var/lib/tftpboot/ignition
COPY node-configs/ .
RUN ct < deep-thought-ignite.yml > deep-thought-ignite.ign \
    && ct < deep-thought-ignite-boot.yml > deep-thought-ignite-boot.ign \
    && ct < allied-mastercomputer-ignite.yml > allied-mastercomputer-ignite.ign \
    && ct < allied-mastercomputer-ignite-boot.yml > allied-mastercomputer-ignite-boot.ign \
    && ct < drknow-ignite.yml > drknow-ignite.ign \
    && ct < drknow-ignite-boot.yml > drknow-ignite-boot.ign

# TODO: I doubt that these are actually still used, and could be an artifact of a difficult
# time getting things to work. Will check on this.
COPY tftpboot/ /var/lib/tftpboot

# Start the webserver on port 8000, which is referenced in the *-ignite-boot.yml files and
# daemonize it to get out of the way of dnsmasq, which we start next. Ready to go!
CMD darkhttpd /var/lib/tftpboot/ignition --port 8000 --daemon \
    && dnsmasq --no-daemon
