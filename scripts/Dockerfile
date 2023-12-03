FROM golang:1.21.4

ENV GO111MODULE=on

ADD . /go/src/github.com/kgoralski/go-crud-template
WORKDIR /go/src/github.com/kgoralski/go-crud-template
#RUN go get github.com/kgoralski/go-crud-template

# RUN go install github.com/swaggo/swag/cmd/swag@latest
# RUN go get -u github.com/swaggo/http-swagger
# RUN go get golang.org/x/crypto/ssh/terminal@v0.16.0

RUN apt-get update && \
    apt-get install -y vim && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN go install cmd/main.go

# RUN swag init -d "./" -g "cmd/main.go"

ENTRYPOINT scripts/go-app-entrypoint.sh

EXPOSE 8080