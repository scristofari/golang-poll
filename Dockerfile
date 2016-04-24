FROM golang:1.6-alpine

RUN apk update && apk add bash git

RUN go get github.com/tools/godep

ADD . /go/src/github.com/golang/scristofari/golang-poll

WORKDIR /go/src/github.com/golang/scristofari/golang-poll

ENV GO15VENDOREXPERIMENT=0
RUN godep get

ENV APP_ENV "dev"
ENV APP_PORT 80
ENV APP_HOST "127.0.0.1"

#RUN go build -o golang-poll .
#ENTRYPOINT ["./golang-poll"]

EXPOSE $APP_PORT


