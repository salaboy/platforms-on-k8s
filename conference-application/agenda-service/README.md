# Agenda Service

## Build from source

To build the application, run the following command:

```shell
go build -o agenda-service cmd/main.go
```

### OpenAPI documentation

We used [oapi-codegen](https://github.com/deepmap/oapi-codegen) to generate all `http.HandlerFunc` from the [openapi.yaml](docs/openapi.yaml) file:

```shell
oapi-codegen -generate chi-server -package main docs/openapi.yaml > api/api.go      
```

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
go test ./... -v
```

To see the code coverage while running tests, execute the following command:

```shell
go test ./...  -coverpkg=./... -coverprofile ./coverage.out
```

## Create Container