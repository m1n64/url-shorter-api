# Infrastructure Service

### Description
General service and main containers

### Stack:
- nginx
- RabbitMQ
- memcached

### How to start
```shell
cp .env.example .env
```
```shell
docker-compose up -d
```

### Ports:
- nginx: 80/443
- RabbitMQ: 15672
- memcached: 11211
