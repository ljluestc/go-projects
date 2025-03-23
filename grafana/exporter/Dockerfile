FROM golang:1.21 AS builder
WORKDIR /app
COPY main.go .
RUN go mod init sample_exporter && go get github.com/prometheus/client_golang/prometheus && go get github.com/prometheus/client_golang/prometheus/promhttp
RUN CGO_ENABLED=0 GOOS=linux go build -o exporter .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/exporter .
EXPOSE 8080
CMD ["./exporter"]