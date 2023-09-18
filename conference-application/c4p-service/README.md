# Agenda Service

## Endpoints

- `/` POST 
- `/` GET
- `/` DELETE
- `/{id}` GET
- `/{id}` DELETE
- `/highlights` GET
- `/day/{day}` GET
- `/health/readiness` GET
- `/health/liveness` GET


## Build & Run from source

```shell
docker-compose up 
```

```shell
go run c4p-service.go
```
## Create Container