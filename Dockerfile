FROM golang:latest

ADD . /go/src/github.com/golang/sparck/golang-poll

RUN go get github.com/tools/godep

WORKDIR /go/src/github.com/golang/sparck/golang-poll
RUN godep get
RUN go run main.go

EXPOSE 9001