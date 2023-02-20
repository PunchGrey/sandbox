# builder image
FROM golang:1.19-alpine as builder

ENV CGO_ENABLED 0
ENV GO111MODULE on
RUN apk --no-cache add git
WORKDIR /go/src/simpleweb
COPY . .
#RUN go test -v ./...
ENV GOARCH amd64
RUN go build -o /bin/simpleweb -v

# final image
FROM alpine:3.14.6
MAINTAINER PunchGrey

RUN apk --no-cache add ca-certificates dumb-init tzdata
COPY --from=builder /bin/simpleweb /bin/simpleweb

USER 65534
ENTRYPOINT ["dumb-init", "--", "/bin/simpleweb"]