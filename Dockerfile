FROM golang:1.7-alpine

MAINTAINER Deniss Gubanov <deniss@gubanov.ee>

WORKDIR /root

COPY ./src/com/github/drewoko /root

ENV PATH $GOPATH/src/github.com/jteeuwen/go-bindata/go-bindata:$PATH

RUN apk add --no-cache git build-base && \
    go get -t github.com/jteeuwen/go-bindata && \
    cd $GOPATH/src/github.com/jteeuwen/go-bindata/go-bindata && \
    go build && \
    cd /root/p3ka && \
    go get -d ./... && \
    go-bindata -o core/bindata.go -pkg core static/* && \
    go build && \
    cp p3ka /bin/p3ka && \
    cd /root && \
    rm -rf p3ka && \
    rm -rf /usr/local/go \
    rm -rf /go  && \
    >application.properties && \
    apk del build-base

ENTRYPOINT ["p3ka"]