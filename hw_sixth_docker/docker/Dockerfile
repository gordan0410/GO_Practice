# syntax=docker/dockerfile:1
FROM golang:1.17-alpine As build-env
WORKDIR /go/src/go_gin
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./app

FROM alpine
WORKDIR /app/
COPY --from=build-env go/src/go_gin/app /app
COPY --from=build-env go/src/go_gin/config.json /app
COPY --from=build-env go/src/go_gin/template /app/template
ENTRYPOINT ./app
