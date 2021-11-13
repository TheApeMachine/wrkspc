FROM rocsys/alpine

# Add a very minimal root filesystem.
ADD alpine-minirootfs-3.14.0-x86_64.tar.gz /

# Setup user account and home path.
ARG USER=guest
ENV HOME /home/$USER

RUN apk add --no-cache ca-certificates sudo git \
	&& echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
	&& chmod 0440 /etc/sudoers.d/$USER \
	&& rm -rf /var/cache/apk/* /usr/share/doc /usr/share/man/ /usr/share/info/* /var/cache/man/* /tmp/* /etc/fstab \
	&& rm -fr /etc/init.d /lib/rc /etc/conf.d /etc/inittab /etc/runlevels /etc/rc.conf \
	&& rm -fr /etc/sysctl* /etc/modprobe.d /etc/modules /etc/mdev.conf /etc/acpi

# Drop to a non-root user account.
USER $USER
WORKDIR $HOME
RUN sudo chown -R $USER:users $HOME

ENTRYPOINT git
