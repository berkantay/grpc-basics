FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal
COPY pkg pkg

RUN go build -o user-management-service ./cmd

FROM alpine:3.14

RUN apk add --no-cache ca-certificates
COPY --from=builder /app/user-management-service /usr/local/bin/

EXPOSE 8080

CMD ["/usr/local/bin/user-management-service"]