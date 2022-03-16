FROM bbsdocker/imageptt:latest


# golang on debian
WORKDIR /opt
RUN curl -L https://dl.google.com/go/go1.16.2.linux-amd64.tar.gz | tar -zxv && \
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
RUN cp 01-config-docker.go.template ptttype/00-config-production.go && \
    mkdir -p /etc/go-pttbbs && cp 01-config.docker.ini /etc/go-pttbbs/production.ini && \
    cp docs/etc/* /etc/go-pttbbs && \
    chown -R bbs .

USER bbs
WORKDIR /srv/go-pttbbs
RUN go build -ldflags "-X github.com/Ptt-official-app/go-pttbbs/types.GIT_VERSION=`git rev-parse --short HEAD` -X github.com/Ptt-official-app/go-pttbbs/types.VERSION=`git describe --tags`" -tags production
RUN ./scripts/rebuild_pttbbs.sh

# mkdir -p /opt/bbs
USER root
RUN mkdir -p /opt/bbs && cp -R /home/bbs/pttbbs /home/bbs/bin /home/bbs/etc /home/bbs/wsproxy /opt/bbs
RUN ./scripts/openrestry.sh

# cmd
WORKDIR /home/bbs
CMD ["sh", "-c", "sudo -iu bbs /home/bbs/bin/shmctl init && sudo -iu bbs /home/bbs/bin/logind && /usr/bin/openresty && sudo -iu bbs /home/bbs/bin/boardd -l 0.0.0.0:5150 && /srv/go-pttbbs/scripts/docker-mand.sh && sudo -iu bbs /srv/go-pttbbs/go-pttbbs -ini production.ini"]

EXPOSE 3456
EXPOSE 5150
EXPOSE 5151
