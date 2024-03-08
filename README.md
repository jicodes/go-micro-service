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

## Add configuarations
- we can use `viper` or `envconfig` package to read configurations from a file.
- But for this project, we will do it manually.

## To-Do
- Add godotenv package to read environment variables from a file.
- Move redisRepo, handlder and model packages to a single `order` package.
- Replace the references to the Repo in Handler with an interface. This will allow us to swap the RedisRepo with any other Repo implementation. For example a PostgresRepo.

```go
type Repo interface {
  Insert(ctx context. Context, order Order) error
  FindByID(ctx context.Context, id uint64) (Order, error)
  DeleteByID(ctx context. Context, id vint64) error
  Update(ctx context. Context, order Order) error
  FindAll(ctx context.Context, page FindAllPage) (FindResult, error)
}
type Order struct {
  Repo Repo
}
```
- Add some tests.