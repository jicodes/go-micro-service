# Microservice with Go, Chi, Redis and Docker

## Install Redis

```sh
brew install redis
```

## Install go-redis

```sh
go get github.com/redis/go-redis/v9
```

## Run Redis with Docker

```sh
docker run -d -p 6379:6379 --name redis redis:alpine
```