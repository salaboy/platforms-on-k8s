package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
	platformFormat "github.com/containerd/containerd/platforms"
)

var platforms = []dagger.Platform{
	"linux/amd64", // a.k.a. x86_64
	"linux/arm64", // a.k.a. aarch64
}

var (
	// the container registry for the multi-platform image
	CONTAINER_REGISTRY      = getEnv("CONTAINER_REGISTRY", "docker.io")
	CONTAINER_REGISTRY_USER = getEnv("CONTAINER_REGISTRY_USER", "salaboy")
)

// util that returns the architecture of the provided platform
func architectureOf(platform dagger.Platform) string {
	return platformFormat.MustParse(string(platform)).Architecture
}

func buildService(ctx context.Context, client *dagger.Client, dir string) ([]*dagger.Container, error) {
	srcDir := client.Host().Directory(dir)

	platformVariants := make([]*dagger.Container, 0, len(platforms))
	for _, platform := range platforms {
		// pull the golang image for the *host platform*. This is
		// accomplished by just not specifying a platform; the default
		// is that of the host.
		ctr := client.Container()
		ctr = ctr.From("golang:1.21-alpine")

		// mount in our source code
		ctr = ctr.WithDirectory("/src", srcDir)
		ctr = ctr.WithMountedCache("/go/pkg/mod", client.CacheVolume("go-mod"))
		ctr = ctr.WithMountedCache("/root/.cache/go-build", client.CacheVolume("go-build"))

		// mount in an empty dir to put the built binary
		ctr = ctr.WithDirectory("/output", client.Directory())

		// ensure the binary will be statically linked and thus executable
		// in the final image
		ctr = ctr.WithEnvVariable("CGO_ENABLED", "0")

		// configure the go compiler to use cross-compilation targeting the
		// desired platform
		ctr = ctr.WithEnvVariable("GOOS", "linux")
		ctr = ctr.WithEnvVariable("GOARCH", architectureOf(platform))

		// build the binary and put the result at the mounted output
		// directory
		ctr = ctr.WithWorkdir("/src")
		ctr = ctr.WithExec([]string{
			"go", "build",
			"-o", "/output/app",
			".",
		})
		// select the output directory
		outputDir := ctr.Directory("/output")

		// wrap the output directory in a new empty container marked
		// with the platform
		binaryCtr := client.
			Container(dagger.ContainerOpts{Platform: platform}).
			WithEntrypoint([]string{"./app"}).
			WithRootfs(outputDir)
		platformVariants = append(platformVariants, binaryCtr)
	}
	return platformVariants, nil
}

func testService(ctx context.Context, client *dagger.Client, dir string) error {
	srcDir := client.Host().Directory(dir)

	// Start Kafka for all services
	kafkaSvc := client.Container().
		From("docker.io/bitnami/kafka:3.4.1-debian-11-r0").
		WithEnvVariable("ALLOW_PLAINTEXT_LISTENER", "yes").
		WithEnvVariable("KAFKA_CFG_LISTENERS", "PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:9094").
		WithEnvVariable("KAFKA_CFG_ADVERTISED_LISTENERS", "PLAINTEXT://kafka:9092,EXTERNAL://kafka:9094").
		WithEnvVariable("KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP", "CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT,PLAINTEXT:PLAINTEXT").
		WithExposedPort(9092)

	client.Container().
		From("docker.io/bitnami/kafka:3.4.1-debian-11-r0").
		WithEntrypoint([]string{"/bin/sh", "-c"}).
		WithExec([]string{
			"kafka-topics.sh --bootstrap-server kafka:9092 --list",
			"echo -e 'Creating kafka topics'",
			"kafka-topics.sh --bootstrap-server kafka:9092 --create --if-not-exists --topic events-topic --replication-factor 1 --partitions 1",
			"echo -e 'Successfully created the following topics:'",
			"kafka-topics.sh --bootstrap-server kafka:9092 --list",
		})

	// accomplished by just not specifying a platform; the default
	// is that of the host.
	ctr := client.Container()

	ctr = ctr.WithEnvVariable("KAFKA_URL", "kafka:9092")

	ctr = ctr.From("golang:1.21-alpine").
		WithServiceBinding("kafka", kafkaSvc)

	if dir == "agenda-service" {
		redisSvc := client.Container().
			From("docker.io/bitnami/redis:7.0.11-debian-11-r12").
			WithEnvVariable("ALLOW_EMPTY_PASSWORD", "yes").
			WithExposedPort(6379)
		ctr = ctr.WithServiceBinding("redis", redisSvc)
		ctr = ctr.WithEnvVariable("REDIS_HOST", "redis")
	}

	if dir == "c4p-service" {
		redisSvc := client.Container().
			From("docker.io/bitnami/redis:7.0.11-debian-11-r12").
			WithEnvVariable("ALLOW_EMPTY_PASSWORD", "yes").
			WithExposedPort(6379)

		postgreSvc := client.Container().
			From("bitnami/postgresql:15.3.0-debian-11-r17").
			WithEnvVariable("POSTGRES_USER", "postgres").
			WithEnvVariable("POSTGRES_PASSWORD", "postgres").
			WithFile("/docker-entrypoint-initdb.d/init.sql", srcDir.File("init.sql")).
			WithExposedPort(5432)
		ctr = ctr.WithServiceBinding("postgres", postgreSvc)
		ctr = ctr.WithEnvVariable("POSTGRES_HOST", "postgres")

		agendaSvc := client.Container().
			From("salaboy/agenda-service-0967b907d9920c99918e2b91b91937b3:v1.0.0").
			WithServiceBinding("kafka", kafkaSvc).
			WithServiceBinding("redis", redisSvc).
			WithEnvVariable("KAFKA_URL", "kafka:9092").
			WithEnvVariable("REDIS_HOST", "redis").
			WithExposedPort(8080)
		ctr = ctr.WithServiceBinding("agenda-service", agendaSvc)
		ctr = ctr.WithEnvVariable("AGENDA_SERVICE_URL", "http://agenda-service:8080")
		notificationsSvc := client.Container().
			From("salaboy/notifications-service-0e27884e01429ab7e350cb5dff61b525:v1.0.0").
			WithServiceBinding("kafka", kafkaSvc).
			WithEnvVariable("KAFKA_URL", "kafka:9092").
			WithExposedPort(8080)
		ctr = ctr.WithServiceBinding("notifications-service", notificationsSvc)
		ctr = ctr.WithEnvVariable("NOTIFICATIONS_SERVICE_URL", "http://notification-service:8080")
	}

	// mount in our source code
	ctr = ctr.WithDirectory("/src", srcDir)
	ctr = ctr.WithMountedCache("/go/pkg/mod", client.CacheVolume("go-mod"))
	ctr = ctr.WithMountedCache("/root/.cache/go-build", client.CacheVolume("go-build"))

	// mount in an empty dir to put the built binary
	ctr = ctr.WithDirectory("/output", client.Directory())

	// ensure the binary will be statically linked and thus executable
	// in the final image
	ctr = ctr.WithEnvVariable("CGO_ENABLED", "0")

	// build the binary and put the result at the mounted output
	// directory
	ctr = ctr.WithWorkdir("/src")
	_, err := ctr.WithExec([]string{
		"go", "test", "-disableTC", "./...",
	}).Sync(ctx)
	return err
}

func publishService(ctx context.Context, client *dagger.Client, dir string, platformVariants []*dagger.Container, tag string) error {
	// publishing the final image uses the same API as single-platform
	// images, but now additionally specify the `PlatformVariants`
	// option with the containers built before.

	imageDigest, err := client.Container().
		Publish(ctx, fmt.Sprintf("%s/%s/%s:%s", CONTAINER_REGISTRY, CONTAINER_REGISTRY_USER, dir, tag), dagger.ContainerPublishOpts{
			PlatformVariants: platformVariants,
		})
	if err != nil {
		fmt.Println("Publishing error: %v ", err)
		return err
	}
	fmt.Println("published multi-platform image with digest", imageDigest)
	return nil
}

func main() {
	var err error
	ctx := context.Background()

	if len(os.Args) < 2 {
		panic(fmt.Errorf("invalid number of arguments: expected command (build, publish-image, helm-package)"))
	}
	client := getDaggerClient(ctx)

	defer client.Close()

	switch os.Args[1] {
	case "build":
		if len(os.Args) < 3 {
			err = fmt.Errorf("invalid number of arguments: expected service path and tag")
			break
		}
		_, err = buildService(ctx, client, os.Args[2])
		if err != nil {
			panic(err)
		}
	case "test":
		err = testService(ctx, client, os.Args[2])
		if err != nil {
			panic(err)
		}
	case "publish":
		pv, err := buildService(ctx, client, os.Args[2])
		if err != nil {
			panic(err)
		}
		err = publishService(ctx, client, os.Args[2], pv, os.Args[3])
		if err != nil {
			panic(err)
		}
	case "all":
		pv, err := buildService(ctx, client, os.Args[2])
		if err != nil {
			panic(err)
		}
		err = testService(ctx, client, os.Args[2])
		if err != nil {
			panic(err)
		}
		err = publishService(ctx, client, os.Args[2], pv, os.Args[3])
		if err != nil {
			panic(err)
		}
	default:
		log.Fatalln("invalid command specified")

	}
}

func getDaggerClient(ctx context.Context) *dagger.Client {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}

	return c
}

// getEnv returns the value of an environment variable, or a fallback value if it is not set.
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
