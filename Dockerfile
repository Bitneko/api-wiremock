## first stage - builder
FROM golang:1.12.5 as development
LABEL maintainer="Ang Ziwei <aireheru@gmail.com>"
WORKDIR /go/src/api-wiremock
RUN go get -v -u golang.org/x/lint/golint

COPY . .

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64 \
  GO111MODULE=on
RUN go build -a -ldflags "-extldflags '-static'" -mod vendor -o ./bin/api-wiremock ./cmd/api-wiremock

## second stage - builder
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=development /go/src/api-wiremock/bin .

EXPOSE 8888

ENTRYPOINT ["./api-wiremock"] 
