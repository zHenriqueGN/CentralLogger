# CentralLogger
An application to centrilize your logs from multiple services

### Prerequisites

* [Docker](https://docs.docker.com/engine/install/)
* [Golang](https://go.dev/doc/install)

### Initial Setup

run postgresql docker-compose

```
docker compose -f deployments/database/postgresql/docker-compose.yaml up -d
```

run rabbit docker-compose

```
docker compose -f deployments/messaging/rabbitmq/docker-compose.yaml up -d
```

Install project dependencies locally

```
go mod download
```

Config your .env file using the example in env.example
```
POSTGRES_HOST=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_PORT=
GRPC_SERVER_PORT=
RABBITMQ_HOST=
RABBITMQ_USER=
RABBITMQ_PASSWORD=
RABBITMQ_PORT=
GRAPHQL_SERVER_PORT=
REST_SERVER_PORT=
```

### Running

Run the tests to guarantee that everything is okay
```
go test -v ./...
```

Initiate your gRPC server

```
go run cmd/grpc_server/main.go
```

Initiate your GraphQL server

```
go run cmd/graphql_server/main.go
```

Initiate your REST server

```
go run cmd/rest_server/main.go
```

## Built With

* [gRPC](https://grpc.io/) - Remote Procedure Call (RPC) framework.
* [Protocol Buffers](https://protobuf.dev/) - Mechanism for serializing structured data.
* [gqlgen](https://gqlgen.com/getting-started/) - Go library for building GraphQL servers
* [RabbitMQ](https://rabbitmq.com/download.html) - Message broker
* [PostgreSQL](https://www.postgresql.org/download/) - Open source relational database
* [go-chi](https://go-chi.io/#/pages/getting_started) - Simple and fast router


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.