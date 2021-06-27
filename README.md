# [API bare bones](#api-bate-bones)
Short answer: Useless API.
Long answer: Well designed over engineered barebones of ANY scalable server side application. It was designed following DDD and clean arch principles protecting its domain logic from everything else. Don't get to much attached to the choice of frameworks and drivers, that's now really the point. They could be easily replaced by any other without having to change a single domain line. (Note that the short answers remains true rs)
- [How to(s)](#make-yourself-at-home)
- [Swagger](#swagger) 
- [TODOs](#todos)

### Make yourself at home
For the following steps [Golang](https://golang.org/doc/install) is gonna be necessary.

##### Run it locally
``go run cmd/api/*``

##### Buid it
``make compile`` (generates binary output at ./build)

##### Run tests
``make test``

##### Run linter
``make metalint``


### Swagger
API docs can be found at [Swagger UI](localhost:3000/docs/v1/library/swagger/index.html).

### TODOs
- Postgres (maybe with sqlc)
- Kafka/RabbitMQ
- Redis