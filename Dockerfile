FROM alpine:3.4

ENV GOPATH=/go PATH=/go/bin:$PATH GODEBUG=netdns=go APP_PORT=80

EXPOSE 80
ENTRYPOINT ["/basicauth-reverseproxy"]

RUN apk add --no-cache ca-certificates go git \
      && mkdir -p /go/src /go/bin \
      && chmod -R 777 /go

ADD . /go/src/github.com/pottava/basicauth-reverseproxy
RUN go build github.com/pottava/basicauth-reverseproxy
