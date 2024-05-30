package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

func dockerclient() *client.Client {
	return must(client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()))
}

func dockerauth() string {
	authConfig := registry.AuthConfig{Username: CFG.DockerUsername, Password: CFG.DockerPassword}
	authBytes := must(json.Marshal(authConfig))
	authBase64 := base64.URLEncoding.EncodeToString(authBytes)
	return authBase64
}

func dockerDeploy(tag string) {
	ctx := context.Background()

	imageName := fmt.Sprintf("%s/%s:%s", CFG.DockerUsername, CFG.DockerRepository, tag)

	cli := dockerclient()
	defer cli.Close()

	log.Println("pulling", imageName)
	out := must(cli.ImagePull(ctx, imageName, image.PullOptions{RegistryAuth: dockerauth()}))
	defer out.Close()
	must(io.Copy(io.Discard, out))
	log.Println("pulled", imageName)

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
		killAll(cli, ctx, containers)
	}

	log.Println("next port is", nextPort)

	config := string(must(os.ReadFile(homefile(applicationConfigFile))))

	cr := must(cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        imageName,
			ExposedPorts: nat.PortSet{"8080/tcp": {}},
			Env:          []string{fmt.Sprintf("CONFIG=%s", config)},
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
		log.Println("killing container:", cont.ID)
		check(cli.ContainerStop(ctx, cont.ID, container.StopOptions{}))
		check(cli.ContainerRemove(ctx, cont.ID, container.RemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}))
	}
}
