# Clothing store brands service

This service is used to make **CRUD** operations on **Brand** model

### How to start:
* Type `go run ./cmd/brands` on your terminal
* Type `make run` on your terminal if you can use Makefile

### How to use:
The service was implemented with gRPC. The default port is 8084,
in order to change port change config.yaml file in
`./pkg/config` package.

#### Methods of the service:
* **CreateBrand**(*Brand*) **returns** (*Brand*)
* **ShowBrand**(*ShowBrandRequest*) **returns** (*Brand*)
* **ListBrand**(*ListBrandRequest*) **returns** (*BrandList*)
* **UpdateBrand**(*UpdateBrandRequest*) **returns** (*Brand*)
* **DeleteBrand**(*DeleteBrandRequest*) **returns** (*Brand*)

In order to get more information on methods open `./cmd/pkg/pb/brands.proto`

### RabbitMQ:
`ShowBrand` method can be called with RabbitMQ. In order to get brand, client has to
publish marshalled request to `brands-queue` queue. Then it will respond with marshalled
result.