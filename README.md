# API bare bones

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/fernandodr19/api-barebones/CI?style=flat-square)](https://github.com/fernandodr19/api-barebones/actions?query=workflow%3AMain)
[![Go Report Card](https://goreportcard.com/badge/github.com/fernandodr19/api-barebones)](https://goreportcard.com/report/github.com/fernandodr19/api-barebones)

Well designed over engineered barebones of ANY server side application. It was designed following DDD and clean arch principles protecting its domain logic from everything else. Don't get to much attached to the choice of frameworks and drivers because that's now really the point. They could be easily replaced by any other without having to change a single domain line.

----------------------------------

- [How to(s)](#make-yourself-at-home)
- [Swagger](#swagger) 
- [Project tree](#project-tree) 
- [TODOs](#todos)

----------------------------------

### Make yourself at home
For the following steps [Golang](https://golang.org/doc/install) is gonna be necessary.

##### Run it locally
``$ go run cmd/api/*``

##### Buid it
``$ make compile`` (generates binary output at ./build)

##### Run tests
``$ make test``

##### Run linter
``$ make metalint``

----------------------------------

### Swagger
Once application is running API docs can be found at [Swagger UI](http://localhost:3000/docs/v1/library/swagger/index.html).

----------------------------------

### Project tree
```bash
$ tree
├── build
│   └── library-api
├── cmd
│   └── api
│       ├── buildcfg.go
│       └── main.go
├── docs
│   └── swagger
│       ├── docs.go
│       ├── swagger.json
│       └── swagger.yaml
├── go.mod
├── go.sum
├── Makefile
├── pkg
│   ├── app.go
│   ├── config
│   │   └── config.go
│   ├── domain
│   │   ├── entities
│   │   │   └── accounts
│   │   │       └── repository.go
│   │   ├── error.go
│   │   └── usecases
│   │       └── accounts
│   │           ├── errors.go
│   │           ├── mocks.gen.go
│   │           └── usecase.go
│   ├── gateway
│   │   ├── api
│   │   │   ├── accounts
│   │   │   │   ├── dosomething.go
│   │   │   │   ├── dosomething_test.go
│   │   │   │   └── handler.go
│   │   │   ├── app.go
│   │   │   ├── middleware
│   │   │   │   └── middleware.go
│   │   │   ├── responses
│   │   │   │   └── responses.go
│   │   │   └── shared
│   │   │       └── shared.go
│   │   └── repositories
│   │       └── accounts.go
│   └── instrumentation
│       └── instrumentation.go
└── README.md
```

----------------------------------

### TODOs
- Postgres (maybe with sqlc)
- Kafka/RabbitMQ
- Redis