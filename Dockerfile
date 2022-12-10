FROM golang:1.17.11

RUN mkdir /workspace
WORKDIR /workspace
RUN apt update && apt upgrade -y

COPY go.mod .
COPY go.sum .
RUN go mod download
