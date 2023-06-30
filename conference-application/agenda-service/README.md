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

### OpenAPI documentation

Execute the application:

```shell
./agenda-service
```

Open your browser and access the following URL:
```http request
http://localhost:8080/openapi/
```

## Test

To execute all tests, run the following command:

```shell
go test ./...
```

To see the code coverage while running tests, execute the following command:

```shell
go test ./...  -coverpkg=./... -coverprofile ./coverage.out
```

## Build from source

## Create Container