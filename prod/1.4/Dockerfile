# Dockerfile

FROM alpine:3.3
MAINTAINER @pottava

LABEL com.github.pottava.application="basicauth-reverseproxy" \
      com.github.pottava.usage="docker run -d -p 8080:80 -e PROXY_URL -e BASIC_AUTH_USER -e BASIC_AUTH_PASS pottava/proxy" \
      com.github.pottava.license="MIT"

ENV GOLANG_VERSION=1.4.3 \
    GOLANG_SRC_URL=https://golang.org/dl/go1.4.3.src.tar.gz \
    GOLANG_SRC_SHA1=486db10dc571a55c8d795365070f66d343458c48 \
    GOPATH=/go \
    PATH=/go/bin:/usr/local/go/bin:$PATH

RUN set -x \
    && apk add --no-cache --virtual build-dependencies \
        bash \
        gcc \
        musl-dev \
        git \
    && apk add --no-cache ca-certificates \
    && wget -q "$GOLANG_SRC_URL" -O golang.tar.gz \
    && echo "$GOLANG_SRC_SHA1  golang.tar.gz" | sha1sum -c - \
    && tar -C /usr/local -xzf golang.tar.gz \
    && rm golang.tar.gz \
    && cd /usr/local/go/src \
    && ./make.bash \
    && mkdir -p "$GOPATH/src" "$GOPATH/bin" \
    && chmod -R 777 "$GOPATH" \
    && go get -u github.com/pottava/basicauth-reverseproxy \
    && apk del --purge build-dependencies \
    && rm -rf /go/pkg /go/src /usr/local/go

EXPOSE 80

ENTRYPOINT ["/go/bin/basicauth-reverseproxy"]
