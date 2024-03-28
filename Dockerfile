FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /marketplace-service ./cmd/api

FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /marketplace-service .

COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./marketplace-service"]
