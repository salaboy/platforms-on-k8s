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

func buildService(ctx context.Context, dir string) ([]*dagger.Container, error) {
	client := getDaggerClient(ctx)

	defer client.Close()

	srcDir := client.Host().Directory(dir)

	platformVariants := make([]*dagger.Container, 0, len(platforms))
	for _, platform := range platforms {
		// pull the golang image for the *host platform*. This is
		// accomplished by just not specifying a platform; the default
		// is that of the host.
		ctr := client.Container()
		ctr = ctr.From("golang:1.20-alpine")

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
	fmt.Println("Artifacts built: %s ", platformVariants)
	return platformVariants, nil
}

func testService(ctx context.Context, dir string) error {
	// run docker compose up and then run go test from inside each service directory
	return nil
}

func publishService(ctx context.Context, dir string, platformVariants []*dagger.Container, tag string) error {
	// publishing the final image uses the same API as single-platform
	// images, but now additionally specify the `PlatformVariants`
	// option with the containers built before.
	client := getDaggerClient(ctx)

	defer client.Close()

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

	switch os.Args[1] {
	case "build":
		if len(os.Args) < 3 {
			err = fmt.Errorf("invalid number of arguments: expected service path and tag")
			break
		}
		_, err = buildService(ctx, os.Args[2])
		if err != nil {
			panic(err)
		}
	case "test":
		err = testService(ctx, os.Args[2])
		if err != nil {
			panic(err)
		}
	case "publish":
		pv, err := buildService(ctx, os.Args[2])
		if err != nil {
			panic(err)
		}
		err = publishService(ctx, os.Args[2], pv, os.Args[3])
		if err != nil {
			panic(err)
		}
	case "all":
		pv, err := buildService(ctx, os.Args[2])
		if err != nil {
			panic(err)
		}
		err = testService(ctx, os.Args[2])
		if err != nil {
			panic(err)
		}
		err = publishService(ctx, os.Args[2], pv, os.Args[3])
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
