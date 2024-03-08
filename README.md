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

## Graceful Shutdown or Graceful Termination

- The application should be able to handle the termination signal and gracefully shutdown the server.
- To handle the termination signal.
  - run the server in `goroutine` so it won't block main thread.
  - make a `chan` to communicate between the main goroutine and the server goroutine.
  - in the main goroutine, listen for the termination signal and send the signal to the server goroutine.
  - `os/signal` package is used to handle the termination signal.
  - `context` package is used to communicate and handle gracefully shutdown of the server.

## Data model

- Install uuid package for generating unique id.

```sh
go get github.com/google/uuid
```

## Add Redis Repository funcs and link with Handlers