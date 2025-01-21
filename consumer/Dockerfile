# Build Stage
FROM golang:1.23.4-alpine3.20 AS builder

WORKDIR /app

COPY . .

RUN go build -o main main.go handler.go

# Run Stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

CMD [ "/app/main" ]
