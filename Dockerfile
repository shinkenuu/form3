FROM golang:1.17.3-alpine

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/shinkenuu/form3

# COPY . .
COPY go.mod ./
COPY main.go ./
COPY client ./client

RUN go get -d -v

RUN go build -o /go/bin/fake_api_client

ENTRYPOINT ["/go/bin/fake_api_client"]