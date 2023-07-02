FROM golang:1.20 as builder

WORKDIR /app

COPY src ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/speedtest-exporter

FROM debian:buster-slim

COPY --from=builder /app/speedtest-exporter /app/speedtest-exporter

EXPOSE 8080

CMD ["/app/speedtest-exporter"]