FROM golang:1.16-alpine as builder

WORKDIR /build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://proxy.golang.org,direct"

RUN apk add --no-cache git

COPY go.mod .
COPY go.sum .

RUN go mod download -x
COPY . .

RUN go build -a -tags 'netgo osusergo' -o /go/bin/foaas-api main.go

FROM alpine
COPY --from=builder go/bin/foaas-api /usr/local/bin

WORKDIR usr/local/bin
ENTRYPOINT [ "foaas-api", "serve" ]
EXPOSE 4000
