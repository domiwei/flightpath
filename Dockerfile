FROM golang:1.18-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go build -o flightpaths ./app/main.go

FROM alpine:3.15
WORKDIR /opt/app
COPY --from=builder /build/flightpaths /opt/app/
EXPOSE 8080
ENTRYPOINT /opt/app/flightpaths
