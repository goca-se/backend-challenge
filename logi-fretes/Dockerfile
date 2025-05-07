FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod .

COPY *.go .

RUN go build -o logifretes-api .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/logifretes-api .

EXPOSE 8080

CMD ["./logifretes-api"] 