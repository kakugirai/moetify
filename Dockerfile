FROM golang:alpine
RUN mkdir /app
COPY . /app
WORKDIR /app
ENV APP_REDIS_ADDR="redis:6379"
RUN apk add git && \
    go build -o app . && \
    adduser -S -D -H -h /app kakugirai
USER kakugirai
CMD ["./app"]