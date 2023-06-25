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

// the container registry for the multi-platform image
const imageRepo = "ttl.sh/marcostest"

// util that returns the architecture of the provided platform
func architectureOf(platform dagger.Platform) string {
	return platformFormat.MustParse(string(platform)).Architecture
}

func buildAndPublishService(ctx context.Context, dir, tag string) {
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

	// publishing the final image uses the same API as single-platform
	// images, but now additionally specify the `PlatformVariants`
	// option with the containers built before.
	imageDigest, err := client.
		Container().
		Publish(ctx, fmt.Sprintf("%s:%s", imageRepo, tag), dagger.ContainerPublishOpts{
			PlatformVariants: platformVariants,
		})
	if err != nil {
		panic(err)
	}
	fmt.Println("published multi-platform image with digest", imageDigest)
}

func main() {
	var err error
	ctx := context.Background()

	if len(os.Args) < 2 {
		panic(fmt.Errorf("invalid number of arguments: expected command (build, publish-image, helm-package)"))
	}

	switch os.Args[1] {
	case "build":
		if len(os.Args) < 4 {
			err = fmt.Errorf("invalid number of arguments: expected service path and tag")
			break
		}
		buildAndPublishService(ctx, os.Args[2], os.Args[3])

	case "helm-publish":
		if len(os.Args) < 3 {
			err = fmt.Errorf("invalid number of arguments: expected chart tag")
			break
		}
		err = helmPublish(ctx, os.Args[2])

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

func helmPublish(ctx context.Context, tag string) error {
	c := getDaggerClient(ctx)
	defer c.Close()

	chartDir := c.Host().Directory("./helm/conference-app")

	helm := c.Container().From("alpine/helm:3.12.1").
		WithMountedDirectory(".", chartDir).
		WithExec([]string{"registry", "login", "-u", "salaboy", "ghcr.io", "--password-stdin"}, dagger.ContainerWithExecOpts{Stdin: os.Getenv("HELM_REGISTRY_PASSWORD")}).
		WithExec([]string{"package", "-u", "."})

	chartOut, err := helm.Stdout(ctx)
	if err != nil {
		return err
	}

	chartPackagePath := strings.TrimSpace(strings.Split(chartOut, ":")[1])

	_, err = helm.WithExec([]string{"push", chartPackagePath, "oci://ghcr.io/salaboy"}).
		ExitCode(ctx)

	return err
}
