# Build the manager binary
FROM golang:1.14-alpine as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY cmd/ cmd/
COPY lib/ lib/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a .

# Use alpine as a base
FROM alpine:3.12
RUN apk add --update --no-cache ca-certificates
WORKDIR /
COPY --from=builder /workspace/obscurus .
COPY public/ public/

ENTRYPOINT ["/obscurus"]