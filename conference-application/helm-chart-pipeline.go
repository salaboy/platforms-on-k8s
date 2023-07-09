package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"dagger.io/dagger"
	platformFormat "github.com/containerd/containerd/platforms"
)

var platforms = []dagger.Platform{
	"linux/amd64", // a.k.a. x86_64
	"linux/arm64", // a.k.a. aarch64
}

var (
	// the container registry for the multi-platform image
	CONTAINER_REGISTRY          = getEnv("CONTAINER_REGISTRY", "docker.io")
	CONTAINER_REGISTRY_USER     = getEnv("CONTAINER_REGISTRY_USER", "salaboy")
	CONTAINER_REGISTRY_PASSWORD = getEnv("CONTAINER_REGISTRY_PASSWORD", "")
)

// util that returns the architecture of the provided platform
func architectureOf(platform dagger.Platform) string {
	return platformFormat.MustParse(string(platform)).Architecture
}

func main() {
	var err error
	ctx := context.Background()

	if len(os.Args) < 2 {
		panic(fmt.Errorf("invalid number of arguments: expected command (build, publish-image, helm-package)"))
	}

	switch os.Args[1] {
	case "package":
		if len(os.Args) < 3 {
			err = fmt.Errorf("invalid number of arguments: expected chart tag")
			break
		}
		chart, err := helmPackage(ctx, os.Args[2])

	case "test":
		if len(os.Args) < 3 {
			err = fmt.Errorf("invalid number of arguments: expected chart tag")
			break
		}
		err = helmTest(ctx, os.Args[2])
	case "publish":
		if len(os.Args) < 3 {
			err = fmt.Errorf("invalid number of arguments: expected chart tag")
			break
		}
		chart, err := helmPackage(ctx, os.Args[2])
		err = helmPublish(ctx, chart)

	case "all":
		chart, err := helmPackage(ctx, os.Args[2])
		if err != nil {
			panic(err)
		}
		err = helmTest(ctx, os.Args[2])
		if err != nil {
			panic(err)
		}
		err = helmPublish(ctx, chart)
		if err != nil {
			panic(err)
		}

	default:
		log.Fatalln("invalid command specified")
	}

	if err != nil {
		panic(err)
	}
}

func getDaggerClient(ctx context.Context) *dagger.Client {
	c, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stderr))
	if err != nil {
		panic(err)
	}

	return c
}

func helmPackage(ctx context.Context, tag string) (string, error) {
	c := getDaggerClient(ctx)
	defer c.Close()

	chartDir := c.Host().Directory("./helm/conference-app")

	helm := c.Container().From("alpine/helm:3.12.1").
		WithMountedDirectory(".", chartDir).
		WithExec([]string{"package", "-u", "."})

	chartOut, err := helm.Stdout(ctx)
	if err != nil {
		return "", err
	}
	return chartOut, nil
}

func helmTest(ctx context.Context, tag string) error {
	// run helm test
	return nil
}

func helmPublish(ctx context.Context, chart string) error {
	c := getDaggerClient(ctx)
	defer c.Close()
	chartPackagePath := strings.TrimSpace(strings.Split(chart, ":")[1])
	helm := c.Container().From("alpine/helm:3.12.1").
		WithExec([]string{"registry", "login", "-u", CONTAINER_REGISTRY_USER, CONTAINER_REGISTRY, "--password-stdin"}, dagger.ContainerWithExecOpts{Stdin: CONTAINER_REGISTRY_PASSWORD}).
		WithExec([]string{"push", chartPackagePath, fmt.Sprintf("%s%s/%s", "oci://", CONTAINER_REGISTRY, CONTAINER_REGISTRY_USER)})
	_, err := helm.Stdout(ctx)

	return err
}

// getEnv returns the value of an environment variable, or a fallback value if it is not set.
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
