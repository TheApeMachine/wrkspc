FROM scratch

ADD alpine-minirootfs-3.15.0-x86_64.tar.gz /

ARG USER=guest
ENV HOME /home/$USER

RUN apk add --no-cache ca-certificates sudo zsh git zsh-autosuggestions \
		&& echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
		&& chmod 0440 /etc/sudoers.d/$USER \
		&& rm -rf /var/cache/apk/* /usr/share/doc /usr/share/man/ /usr/share/info/* /var/cache/man/* /tmp/* /etc/fstab \
		&& rm -fr /etc/init.d /lib/rc /etc/conf.d /etc/inittab /etc/runlevels /etc/rc.conf \
		&& rm -fr /etc/sysctl* /etc/modprobe.d /etc/modules /etc/mdev.conf /etc/acpi

WORKDIR $HOME

RUN sh -c "$(wget https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh -O -)" \
		&& echo "source /usr/share/zsh/plugins/zsh-autosuggestions/zsh-autosuggestions.zsh" >> ~/.zshrc

RUN sudo chown -R $USER:users $HOME
USER $USER

ENTRYPOINT /bin/zsh
