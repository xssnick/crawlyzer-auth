FROM golang:1.13 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service

###################### Building lightweight

FROM alpine:latest

WORKDIR /app

COPY ./migrations /app/migrations
COPY --from=builder /app/service .

EXPOSE 8080

ENTRYPOINT ["./service"]