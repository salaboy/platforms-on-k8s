module github.com/salaboy/platforms-on-k8s/conference-application/c4p-service

go 1.21

toolchain go1.21.0

require (
	github.com/dapr/go-sdk v1.8.0
	github.com/deepmap/oapi-codegen v1.13.0
	github.com/go-chi/chi v1.5.4
	github.com/go-chi/chi/v5 v5.0.10
	github.com/google/uuid v1.3.0
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.8.4
	github.com/testcontainers/testcontainers-go v0.23.0
	github.com/testcontainers/testcontainers-go/modules/compose v0.23.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/genproto v0.0.0-20230410155749-daa745c078e1 // indirect
	google.golang.org/grpc v1.55.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/cucumber/godog => github.com/laurazard/godog v0.0.0-20220922095256-4c4b17abdae7
