FROM golang:1.22-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 go build -o main ./cmd/semver/main.go

FROM scratch

COPY --from=builder /usr/src/app/main /main

ENTRYPOINT ["/main"]
