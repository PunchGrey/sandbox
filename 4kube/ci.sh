#!/bin/bash
docker run -id --name=4kube_build -v /home/punchgrey/progect/go/src/sandbox/4kube:/go/src/4kube  golang:alpine
docker exec -it 4kube_build go build -o /go/src/4kube/4kube /go/src/4kube/main.go
docker rm -f 4kube_build

docker build -t punchgrey/4kube:1 .