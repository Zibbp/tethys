FROM golang:1.19 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN CGO_ENABLED=1 GOOS=linux go build -o tethys cmd/sever/main.go

FROM alpine:latest AS production

COPY --from=builder /app .

EXPOSE 28542

CMD ["./tethys"]