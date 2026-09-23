FROM golang:1.26.0-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /worker ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /worker .
ENTRYPOINT ["./worker"]