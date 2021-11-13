FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && apt install --no-install-recommends -y bash-completion build-essential \
    cmake libcurl4 libcurl4-openssl-dev libssl-dev libxml2 libxml2-dev libssl1.1 \
    pkg-config ca-certificates xclip git wget golang-go

WORKDIR /tmp/src/
RUN git clone https://github.com/lastpass/lastpass-cli.git

WORKDIR /tmp/src/lastpass-cli
RUN make && make install && rm -rf /tmp/src

WORKDIR /root
COPY main.go .

# Why bother compiling...
ENTRYPOINT ["go", "run", "main.go"]