[![build status](https://img.shields.io/github/workflow/status/fernandodr19/api-barebones/Code%20Quality)](https://github.com/fernandodr19/api-barebones/actions)

![fluxo de trabalho de exemplo](https://github.com/github/docs/actions/workflows/main.yml/badge.svg)
![fluxo de trabalho de exemplo](https://github.com/fernandodr19/api-barebones/.github/workflows/main.yml/badge.svg)
![fluxo de trabalho de exemplo](https://github.com/fernandodr19/api-barebones/workflows/main.yml/badge.svg)

<p align="left">
  <a href="https://github.com/actions/setup-go/actions"><img alt="GitHub Actions status" src="https://github.com/actions/setup-go/workflows/build-test/badge.svg"></a>
</p>

![GitHub tag](https://img.shields.io/github/tag-date/fernandodr19/api-barebones.svg?style=popout)

[![build status](https://img.shields.io/github/workflow/status/fernandodr19/api-barebones/Code%20Quality)](https://github.com/fernandodr19/api-barebones/actions)

[![codecov](https://codecov.io/gh/fernandodr19/api-barebones/branch/master/graph/badge.svg)](https://codecov.io/gh/fernandodr19/api-barebones) 

# API bare bones
Short answer: Useless API.

Long answer: Well designed over engineered barebones of ANY server side application. It was designed following DDD and clean arch principles protecting its domain logic from everything else. Don't get to much attached to the choice of frameworks and drivers because that's now really the point. They could be easily replaced by any other without having to change a single domain line. (Note that the short answers remains true rs)

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