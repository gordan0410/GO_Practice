FROM golang:1.17

WORKDIR /go/src/
COPY . .
RUN cd ./go_gin ; go mod tidy
RUN apt-get update && apt-get install -y \
    vim