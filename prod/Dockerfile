FROM alpine:3.10

ENV APP_PORT=80

RUN apk add --no-cache ca-certificates

RUN apk --no-cache add --virtual build-dependencies gcc musl-dev go git \
    && export GOPATH=/go \
    && export PATH=$GOPATH/bin:$PATH \
    && mkdir $GOPATH \
    && chmod -R 777 $GOPATH \
    && go get -u github.com/pottava/basicauth-reverseproxy \
    && cd /go/src/github.com/pottava/basicauth-reverseproxy \
    && go build -ldflags "-s -w" \
    && mv basicauth-reverseproxy /usr/bin \
    && apk del --purge -r build-dependencies \
    && rm -rf /go

EXPOSE 80

CMD ["basicauth-reverseproxy"]
