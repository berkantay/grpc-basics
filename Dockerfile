FROM golang:1.19-alpine AS builder

WORKDIR /app

RUN apk add build-base
RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev librdkafka-dev pkgconf musl

COPY go.mod go.sum ./
RUN go mod download


COPY cmd cmd
COPY broker broker
COPY database database
COPY grpc grpc
COPY user user
COPY model model

RUN go build -tags musl -ldflags="-X 'main.Version=v1.0.0'" -o user-management-service ./cmd

FROM alpine:3.14

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/user-management-service /usr/local/bin/

EXPOSE 8080

CMD ["/usr/local/bin/user-management-service"]