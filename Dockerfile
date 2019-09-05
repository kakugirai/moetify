# Base
FROM golang:latest as builder
LABEL maintainer="Girai Kaku <kakugirai@gmail.com>"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
ENV APP_REDIS_ADDR="redis:6379"
EXPOSE 8080
CMD ["./main"]