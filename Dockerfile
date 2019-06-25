FROM golang:1.12.5 as development
LABEL maintainer="Ang Ziwei <aireheru@gmail.com>"
WORKDIR /go/src/api-wiremock
RUN go get -v -u golang.org/x/lint/golint

COPY . .

# Download dependencies
ENV GO111MODULE=on
RUN go mod vendor

# Install the package
RUN go install -v ./...

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/api-wiremock .

# Run the binary program produced by `go install`
CMD ["api-wiremock"]

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=development /go/bin/api-wiremock .

EXPOSE 8888

CMD ["./api-wiremock"] 
