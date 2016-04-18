FROM golang:latest

RUN go get github.com/tools/godep

ADD . /go/src/github.com/golang/scristofari/golang-poll
WORKDIR /go/src/github.com/golang/scristofari/golang-poll
RUN godep get

EXPOSE 80