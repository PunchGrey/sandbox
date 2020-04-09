#!/bin/bash
docker run -it --rm -v /home/punchgrey/progect/go/src/sandbox/4kube:/go/src/4kube  golang:alpine sh
go build -o /go/src/4kube/4kube /go/src/4kube/main.go
exit
