# moetify

[![Actions Status](https://github.com/kakugirai/moetify/workflows/Go/badge.svg)](https://github.com/kakugirai/moetify/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/kakugirai/moetify)](https://goreportcard.com/report/github.com/kakugirai/moetify)

A minimal URL shortening service.

- [gorilla/mux](https://github.com/gorilla/mux) for routing.
- [urfave/negroni](https://github.com/urfave/negroni) for middleware.
- [go-redis/redis](https://github.com/go-redis/redis) for caching.

## Build moetify locally

```bash
docker build -t moetify:latest .
```

## Deploy on Kubernetes

Deploy redis and moetify.

```bash
kubectl apply -f k8s/redis-deployment.yml
```

```bash
kubectl apply -f k8s/moetify-deployment.yml
```

Create a redis service.

```bash
kubectl apply -f redis-service.yml
```

Create a moetify service of type LoadBalancer.

```
kubectl expose deployment moetify-app --name moetify-lb-service \
    --type LoadBalancer --port 80 --target-port 8080
```


## API

Create a tiny url.

```bash
curl -X POST http://moetify.vim.moe/api/shorten -d '{
    "url": "http://example.com/",
    "expiration_in_minutes": 5
}'
```

Get more tiny url information.

```
curl -X GET 'http://moetify.vim.moe/api/info?shortlink=B22ul6TZwiRl'
```

[http://moetify.vim.moe/B33spMi94CO4](http://moetify.vim.moe/B22ul6TZwiRl) will be redirected to [http://example.com/](http://example.com/).
