FROM golang:latest as build

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
ENV CGO_ENABLED 0
ENV GO111MODULE on

RUN apt update -y && apt upgrade -y ca-certificates
RUN mkdir -p /go/{src,bin,pkg}

ADD . /go/src/github.com/takaishi/k8s-github-auth
WORKDIR /go/src/github.com/takaishi/k8s-github-auth
RUN go get
RUN go build

FROM alpine:latest as app
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=build /go/src/github.com/takaishi/k8s-github-auth/k8s-github-auth /k8s-github-auth

ENTRYPOINT ["/k8s-github-auth"]
