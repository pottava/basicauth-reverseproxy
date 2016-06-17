# Dockerfile

FROM alpine:3.4
MAINTAINER @pottava

LABEL com.github.pottava.application="basicauth-reverseproxy" \
      com.github.pottava.usage="docker run -d -p 8080:80 -e PROXY_URL -e BASIC_AUTH_USER -e BASIC_AUTH_PASS pottava/proxy" \
      com.github.pottava.license="MIT"

ENV GOPATH=/go PATH=/go/bin:$PATH GODEBUG=netdns=go APP_PORT=80

RUN apk add --no-cache ca-certificates \
      && apk --no-cache add --virtual build-dependencies go git \
      && mkdir -p /go/src /go/bin \
      && chmod -R 777 /go \
      && go get -u github.com/pottava/basicauth-reverseproxy \
      && apk del --purge build-dependencies \
      && rm -rf /go/pkg /go/src

EXPOSE 80

ENTRYPOINT ["/go/bin/basicauth-reverseproxy"]
