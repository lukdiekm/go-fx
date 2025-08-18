FROM golang:1.24.5 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-fx

FROM docker
WORKDIR /root/
COPY --from=builder /go-fx .
RUN chmod +x /go-fx
CMD ["/go-fx"]