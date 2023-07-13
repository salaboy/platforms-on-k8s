package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"dagger.io/dagger"
)

var (
	// the container registry for the multi-platform image
	CONTAINER_REGISTRY          = getEnv("CONTAINER_REGISTRY", "docker.io")
	CONTAINER_REGISTRY_USER     = getEnv("CONTAINER_REGISTRY_USER", "salaboy")
	CONTAINER_REGISTRY_PASSWORD = getEnv("CONTAINER_REGISTRY_PASSWORD", "")
)

func main() {
	var err error
	ctx := context.Background()

	if len(os.Args) < 2 {
		panic(fmt.Errorf("invalid number of arguments: expected command (package, test, publish, all)"))
	}

	client := getDaggerClient(ctx)
	defer client.Close()

	switch os.Args[1] {
	case "package":
		_, _, err := helmPackage(ctx, client)
		if err != nil {
			fmt.Println("Packaging error: %v ", err)
		}
	case "test":
		err = helmTest(ctx)
	case "publish":
		hc, chart, err := helmPackage(ctx, client)
		if err != nil {
			fmt.Println("Packaging error: %v ", err)
		}
		err = helmPublish(ctx, hc, chart)
		if err != nil {
			fmt.Println("Publishing error: %v ", err)
		}
	case "all":
		hc, chart, err := helmPackage(ctx, client)
		if err != nil {
			panic(err)
		}
		err = helmTest(ctx)
		if err != nil {
			panic(err)
		}
		err = helmPublish(ctx, hc, chart)
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

func helmPackage(ctx context.Context, c *dagger.Client) (*dagger.Container, string, error) {
	chartDir := c.Host().Directory("./helm/conference-app")

	helm := c.Container().From("alpine/helm:3.12.1").
		WithMountedDirectory(".", chartDir).
		WithExec([]string{"repo", "add", "bitnami", "https://charts.bitnami.com/bitnami"}).
		WithExec([]string{"dependency", "build"}).
		WithExec([]string{"package", "-u", "."})

	chartOut, err := helm.Stdout(ctx)
	if err != nil {
		return nil, "", err
	}
	return helm, chartOut, nil
}

func helmTest(ctx context.Context) error {
	// run helm test
	return nil
}

func helmPublish(ctx context.Context, c *dagger.Container, chart string) error {
	chartPackagePath := strings.TrimSpace(strings.Split(chart, ":")[1])
	helm := c.WithExec([]string{"registry", "login", "-u", CONTAINER_REGISTRY_USER, CONTAINER_REGISTRY, "--password-stdin"}, dagger.ContainerWithExecOpts{Stdin: CONTAINER_REGISTRY_PASSWORD}).
		WithExec([]string{"push", chartPackagePath, fmt.Sprintf("%s%s/%s", "oci://", CONTAINER_REGISTRY, CONTAINER_REGISTRY_USER)})
	out, err := helm.Stdout(ctx)
	fmt.Sprintln("Publish out: %s ", out)
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
