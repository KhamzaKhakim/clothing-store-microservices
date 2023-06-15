# Clothing store clothes service

This service is used to make **CRUD** operations on **Clothe** model

### How to start:
* Type `go run ./cmd/clothes` on your terminal
* Type `make run` on your terminal if you can use Makefile

### How to use:
The service was implemented with gRPC. The default port is 8083,
in order to change port change config.yaml file in
`./pkg/config` package.

#### Methods of the service:
* **CreateClothe**(*Clothe*) **returns** (*Clothe*)
* **ShowClothe**(*ShowClotheRequest*) **returns** (*Clothe*)
* **ListClothe**(*ListClotheRequest*) **returns** (*ClotheList*)
* **UpdateClothe**(*UpdateClotheRequest*) **returns** (*Clothe*)
* **DeleteClothe**(*DeleteClotheRequest*) **returns** (*Clothe*) {}

In order to get more information on methods open `./cmd/pkg/pb/clothes.proto`

### RabbitMQ:
`ShowClothe` method can be called with RabbitMQ. In order to get clothe, client has to
publish marshalled request to `clothes-queue` queue. Then it will respond with marshalled
result.