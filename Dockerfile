FROM golang:1.8-alpine

MAINTAINER Deniss Gubanov <deniss@gubanov.ee>

WORKDIR /root

COPY ./src/com/github/drewoko /root

ENV PATH $GOPATH/src/github.com/jteeuwen/go-bindata/go-bindata:$PATH

RUN apk add --no-cache git build-base nodejs && \
    go get -t github.com/jteeuwen/go-bindata && \
    cd $GOPATH/src/github.com/jteeuwen/go-bindata/go-bindata && \
    go build && \
    cd /root/p3ka && \
    go get -d ./... && \
    cd static && \
    npm install && \
    npm install -g webpack && \
    npm run build:prod && \
    cd .. && \
    go-bindata -o core/bindata.go -pkg core static/dist/* && \
    go build && \
    cp p3ka /bin/p3ka && \
    cd /root && \
    rm -rf p3ka && \
    rm -rf /usr/local/go \
    rm -rf /go  && \
    >application.properties && \
    npm uninstall webpack  && \
    apk del build-base git nodejs

ENTRYPOINT ["p3ka"]