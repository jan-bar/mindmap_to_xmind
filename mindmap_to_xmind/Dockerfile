FROM ubuntu:22.04

RUN apt-get update
RUN apt-get install -y jq git curl make ca-certificates openssh-client wget \
    build-essential gcc-mingw-w64 g++-mingw-w64 gcc-mingw-w64-i686 g++-mingw-w64-i686 \
    file --no-install-recommends

ENV GOPROXY https://goproxy.cn,direct
ENV GO111MODULE on
ENV G_EXPERIMENTAL=true
ENV G_MIRROR=https://golang.google.cn/dl/
ENV G_HOME=/root/g
ENV GOPATH $G_HOME/path
ENV GOROOT $G_HOME/go
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH

RUN curl -sL $(curl -s https://api.github.com/repos/voidint/g/releases/latest | jq -r ".assets[].browser_download_url | select(match(\"linux-amd64\"))") | tar -xzC /sbin/ && \
    g install $(g ls-remote stable | tail -n1)

ENTRYPOINT ["/root/xmind/build.sh"]
