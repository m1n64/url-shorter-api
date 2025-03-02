## URL SHORTER API | WIP (in progress)

### Description
This project is a simple URL shorter API, written on Go in microservice architecture. It is a simple REST API with gRPC services.

### Links:
- [github](https://github.com/m1n64/url-shorter-api)

### Checklist:
- [x] user-service (Authorization, token validation, user info) (gRPC)
- [x] link-service (URL shorter information) (gRPC, RabbitMQ)
- [ ] analytics-service (URL redirect analytics) (gRPC, RabbitMQ)
- [x] infrastructure-service (OpenResty, memcache, rabbitmq)
- [ ] gateway-service (REST API, maybe Krakend)
- [ ] API docs service

### Peculiarities:
- Microservice architecture
- Authorization and authentication
- Transactional database
- Relational, key-value, column databases
- gRPC, RabbitMQ, REST API

### Startup (makefile)
Init (for create network):
```shell
make network
```
```shell
make up
```

### Startup (docker-compose):
```shell
docker network create tidy-url-network
```
```shell
docker-compose -f infrastructure-service/docker-compose.yml up -d
```
```shell
cd user-service
```
```shell
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d
```
and all the services will be started. (yep, in future i will add k8s support, because I know that this is a "crutch" solution).
