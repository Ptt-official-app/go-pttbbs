FROM bbsdocker/imageptt:latest

# golang on debian
# THE_ARCH: uname -m
# THE_GOARCH: amd64 or arm64
WORKDIR /opt
RUN THE_ARCH=`uname -m` && \
    if [ "${THE_ARCH}" = "aarch64" -o "${THE_ARCH}" = "arm64"  ]; then THE_GOARCH=arm64; else THE_GOARCH=amd64; fi && \
    echo "THE_ARCH: ${THE_ARCH} THE_GOARCH: ${THE_GOARCH}" && \
    sleep 3 && \
    curl -L https://dl.google.com/go/go1.24.1.linux-${THE_GOARCH}.tar.gz | tar -zxv && \
    mv go /usr/local

RUN DEBIAN_FRONTEND=noninteractive &&  \
    apt update && apt install -y \
        bmake \
        gcc \
        g++ \
        libc6-dev \
        libevent-dev \
        pkg-config \
        ccache \
        clang \
        libgrpc++-dev \
        protobuf-compiler \
        protobuf-compiler-grpc \
        libgflags-dev

ENV GOROOT=/usr/local/go
ENV PATH=${PATH}:/usr/local/go/bin:/home/bbs/bin:/opt/bbs/bin

# go-pttbbs
COPY . /srv/go-pttbbs

WORKDIR /srv/go-pttbbs
RUN cp docs/config/config-docker.go ptttype/01-config-docker.go && \
    mkdir -p /etc/go-pttbbs && cp docs/config/01-config.docker.ini /etc/go-pttbbs/production.ini && \
    cp docs/etc/* /etc/go-pttbbs && \
    chown -R bbs .

USER bbs
WORKDIR /srv/go-pttbbs
RUN cd go-pttbbs; go build -ldflags "-X github.com/Ptt-official-app/go-pttbbs/types.GIT_VERSION=`git rev-parse --short HEAD` -X github.com/Ptt-official-app/go-pttbbs/types.VERSION=`git describe --tags`" -tags docker; cd ..
# XXX all-guest: no need to rebuild pttbbs
# RUN ./scripts/rebuild_pttbbs.sh

# XXX all-guest: no need to run openresty
# mkdir -p /opt/bbs
USER root
# RUN mkdir -p /opt/bbs && cp -R /home/bbs/pttbbs /home/bbs/bin /home/bbs/etc /home/bbs/wsproxy /opt/bbs
# RUN ./scripts/openrestry.sh

# XXX all-guest: no need to run shmctl, logind, openresty, boardd, mand.
# cmd
WORKDIR /home/bbs
CMD ["sh", "-c", "sudo -iu bbs /srv/go-pttbbs/go-pttbbs/go-pttbbs -ini production.ini"]

EXPOSE 3456
# XXX all-guest: no need to expose boardd and mand.
# EXPOSE 5150
# EXPOSE 5151
