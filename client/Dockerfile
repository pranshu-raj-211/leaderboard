FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY lb_client/ ./lb_client/

RUN go build -o main ./lb_client/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

ENTRYPOINT ["./main"]