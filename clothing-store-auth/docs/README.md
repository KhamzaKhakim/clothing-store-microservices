# Clothing store auth service

This service is used to make **CRUD** operations on **User** and **Token** models

### How to start:
* Type `go run ./cmd/auth` on your terminal
* Type `make run` on your terminal if you can use Makefile

### How to use:
The service was implemented with gRPC. The default port is 8082,
in order to change port change config.yaml file in
`./pkg/config` package.

#### Methods of the service:
* **Register**(*RegisterRequest*) **returns** (*RegisterResponse*)
* **Login**(*LoginRequest*) **returns** (*LoginResponse*)
* **Activate**(*ActivateRequest*) **returns** (*ActivateResponse*)
* **Authenticate**(*AuthenticateRequest*) **returns** (*AuthenticateResponse*)
* **Authorize**(*AuthorizeRequest*) **returns** (*AuthorizeResponse*)
* **DeleteUser**(*DeleteUserRequest*) **returns** (*DeleteUserResponse*)

In order to get more information on methods open `./cmd/pkg/pb/auth.proto`

### RabbitMQ:
`Login` method can be called with RabbitMQ. In order to get token, client has to
publish marshalled request to `auth-queue` queue. Then it will respond with marshalled
result.