FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /goflow cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /goflow .
COPY examples/ ./examples/
ENTRYPOINT ["./goflow"]
CMD ["run", "-file", "examples/order_process.json"] 