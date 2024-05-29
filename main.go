package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

var (
	username    = flag.String("username", "", "dockerhub username")
	password    = flag.String("password", "", "dockerhub password or token")
	repository  = flag.String("repository", "", "dockerhub repository")
	tag         = flag.String("tag", "", "image tag")
	isDebugMode = flag.Bool("debug", false, "debug mode")
)

func main() {
	flag.Parse()
	ctx := context.Background()

	if !*isDebugMode {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
				os.Exit(1)
			}
		}()
	}

	if *username == "" || *password == "" || *repository == "" || *tag == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	authConfig := registry.AuthConfig{Username: *username, Password: *password}
	authBytes := must(json.Marshal(authConfig))
	authBase64 := base64.URLEncoding.EncodeToString(authBytes)

	imageName := fmt.Sprintf("%s/%s:%s", *username, *repository, *tag)
	log.Println("deploying", imageName)

	cli := must(client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()))
	defer cli.Close()

	log.Println("pulling", imageName)
	out := must(cli.ImagePull(ctx, imageName, image.PullOptions{RegistryAuth: authBase64}))
	defer out.Close()
	must(io.Copy(io.Discard, out))
	log.Println("✅ pulled", imageName)

	containers := must(cli.ContainerList(ctx, container.ListOptions{All: true}))

	var nextPort string = "8080"
	var shouldKillOneContainer bool
	var containerToKill types.Container

	if len(containers) == 1 {
		cont := containers[0]
		for _, port := range cont.Ports {
			if port.PublicPort == 8080 || port.PublicPort == 8081 {
				containerToKill = cont
				shouldKillOneContainer = true
				if port.PublicPort == 8080 {
					nextPort = "8081"
				} else {
					nextPort = "8080"
				}
			}
		}
	} else if len(containers) > 1 {
		kill(cli, ctx, containers[0])
	}

	cr := must(cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        imageName,
			ExposedPorts: nat.PortSet{"8080/tcp": {}},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"8080/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: nextPort}},
			},
		},
		&network.NetworkingConfig{},
		v1.DescriptorEmptyJSON.Platform,
		"",
	))
	log.Println("created container:", cr.ID)

	check(cli.ContainerStart(ctx, cr.ID, container.StartOptions{}))
	log.Println("started container", cr.ID)

	if shouldKillOneContainer {
		killAll(cli, ctx, []types.Container{containerToKill})
	}
}

func killAll(cli *client.Client, ctx context.Context, containers []types.Container) {
	for _, cont := range containers {
		kill(cli, ctx, cont)
	}
}

func kill(cli *client.Client, ctx context.Context, cont types.Container) {
	log.Println("killing container:", cont.ID)
	check(cli.ContainerStop(ctx, cont.ID, container.StopOptions{}))
	check(cli.ContainerRemove(ctx, cont.ID, container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func must[T any](t T, err error) T {
	check(err)
	return t
}
