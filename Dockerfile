FROM --platform=linux/amd64 bitnami/minideb:latest AS builder

ARG USERNAME=n00b
ARG USER_UID=1000
ARG USER_GID=$USER_UID

ENV GOVERSION="1.19.3"
ENV PATH /home/$USERNAME/.local/bin:/usr/bin:/usr/local/bin:/usr/local/go/bin:$PATH

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && install_packages sudo ca-certificates build-essential ssh git-core wget \
    && rm -rf /usr/local/go \
    && wget "https://golang.org/dl/go${GOVERSION}.linux-amd64.tar.gz" -4 \
    && tar -C /usr/local -xvf "go${GOVERSION}.linux-amd64.tar.gz" \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME

WORKDIR /home/$USERNAME
COPY . .
RUN chown -R $USERNAME:$USERNAME /home/$USERNAME
USER $USERNAME

RUN go build -o wrkspc

FROM --platform=linux/amd64 gcr.io/distroless/base-debian11 AS runtime
WORKDIR /tmp/wrkspc

COPY --from=builder /tmp/wrkspc/wrkspc /tmp/wrkspc/
COPY --from=builder /tmp/wrkspc/docs /tmp/wrkspc/

CMD [ "./wrkspc" ]
