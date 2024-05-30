FROM golang:1.22.3-bookworm AS builder

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/github.com/frisby/go-duckdb-distroless

COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go

RUN go mod download && \
    go mod verify && \
    CGO_ENABLED=1 go build -o /go/bin/go-duckdb-distroless main.go

# Almost identical to gcr.io/distroless/base but with a few added c/c++ libraries:
# - libgomp1
# - libstdc++6
# - libgcc-s1
FROM gcr.io/distroless/cc

COPY --from=builder /go/bin/go-duckdb-distroless /
CMD ["/go-duckdb-distroless"]
