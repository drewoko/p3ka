#!/bin/bash
git pull
go build .
docker build --no-cache -t drewoko/p3ka .
docker stop p3ka
docker rm p3ka
docker run -d --restart=on-failure -e 'GIN_MODE=release' -v '/etc/ssl/certs/ca-certificates.crt:/etc/ssl/certs/ca-certificates.crt' -v /srv/p3ka/p3ka.db:/root/p3ka.db -v /srv/p3ka/application.properties:/root/application.properties -p 127.0.0.1:8087:8087 --name p3ka drewoko/p3ka