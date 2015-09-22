FROM golang:latest

RUN go get github.com/tools/godep
RUN go get github.com/codegangsta/gin

ADD . /go/src/github.com/golang/sparck/golang-poll
WORKDIR /go/src/github.com/golang/sparck/golang-poll
RUN godep get

EXPOSE 9001