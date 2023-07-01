FROM golang:1.20

WORKDIR /app

COPY src ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /speedtest-exporter

EXPOSE 8080

CMD ["/speedtest-exporter"]
