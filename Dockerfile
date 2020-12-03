FROM bbsdocker/imageptt:letflatcenable


# golang on debian
WORKDIR /opt
RUN curl -L https://dl.google.com/go/go1.15.3.linux-amd64.tar.gz | tar -zxv && \
    mv go /usr/local

RUN DEBIAN_FRONTEND=noninteractive &&  \
    apt install -y \
        bmake \
        gcc \
        g++ \
        libc6-dev \
        libevent-dev \
        pkg-config \
        ccache \
        clang

ENV GOROOT=/usr/local/go
ENV PATH=${PATH}:/usr/local/go/bin:/home/bbs/bin:/opt/bbs/bin

# go-bbs
COPY . /srv/go-bbs

WORKDIR /srv/go-bbs
RUN cp 01-config-docker.go.template ptttype/00-config-production.go && \
    mkdir -p /etc/go-bbs && cp 01-config.docker.ini /etc/go-bbs/production.ini && \
    chown -R bbs .

USER bbs
WORKDIR /srv/go-bbs/go-bbs
RUN go build -tags production
WORKDIR /srv/go-bbs
RUN ./scripts/rebuild_pttbbs.sh

# mkdir -p /opt/bbs
USER root
RUN mkdir -p /opt/bbs && cp -R /home/bbs/bin /home/bbs/etc /home/bbs/wsproxy /opt/bbs

# cmd
WORKDIR /home/bbs
CMD ["sh", "-c", "sudo -iu bbs /home/bbs/bin/shmctl init && sudo -iu bbs /home/bbs/bin/logind && /usr/bin/openresty && sudo -iu bbs /srv/go-bbs/go-bbs/go-bbs -ini production.ini"]

EXPOSE 3456
