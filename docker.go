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

var (
	ctx          = context.Background()
	dockerclient = must(client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()))
)

func dockerauth() string {
	authConfig := registry.AuthConfig{Username: CFG.DockerUsername, Password: CFG.DockerPassword}
	authBytes := must(json.Marshal(authConfig))
	authBase64 := base64.URLEncoding.EncodeToString(authBytes)
	return authBase64
}

func formatImage(tag string) string {
	return fmt.Sprintf("%s/%s:%s", CFG.DockerUsername, CFG.DockerRepository, tag)
}

func readConfig() string {
	return string(must(os.ReadFile(homefile(applicationConfigFile))))
}

func pullImage(imageName string) {
	log.Println("pulling", imageName)
	out := must(dockerclient.ImagePull(ctx, imageName, image.PullOptions{RegistryAuth: dockerauth()}))
	defer out.Close()
	must(io.Copy(io.Discard, out))
	log.Println("pulled", imageName)
}

func list() []types.Container {
	return must(dockerclient.ContainerList(ctx, container.ListOptions{All: true}))
}

func dockerDeploy(tag string) {
	imageName := formatImage(tag)
	pullImage(imageName)

	containers := list()

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
		killAll(containers)
	}

	log.Println("next port is", nextPort)

	cr := must(dockerclient.ContainerCreate(
		ctx,
		&container.Config{
			Image:        imageName,
			ExposedPorts: nat.PortSet{"8080/tcp": {}},
			Env:          []string{fmt.Sprintf("CONFIG=%s", readConfig())},
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

	check(dockerclient.ContainerStart(ctx, cr.ID, container.StartOptions{}))
	log.Println("started container", cr.ID)

	if shouldKillOneContainer {
		killAll([]types.Container{containerToKill})
	}
}

func killAll(containers []types.Container) {
	for _, cont := range containers {
		log.Println("killing container:", cont.ID)
		check(dockerclient.ContainerStop(ctx, cont.ID, container.StopOptions{}))
		check(dockerclient.ContainerRemove(ctx, cont.ID, container.RemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}))
	}
}

func runJob(job string) {
	imageName := formatImage(job)
	log.Println("running job:", imageName)
	pullImage(imageName)
	cr := must(dockerclient.ContainerCreate(
		ctx,
		&container.Config{
			Image: imageName,
			Env:   []string{fmt.Sprintf("CONFIG=%s", readConfig())},
		},
		&container.HostConfig{},
		&network.NetworkingConfig{},
		v1.DescriptorEmptyJSON.Platform,
		"",
	))
	log.Println("created container", cr.ID)
	check(dockerclient.ContainerStart(ctx, cr.ID, container.StartOptions{}))
	log.Println("started", cr.ID)

	out := must(dockerclient.ContainerLogs(ctx, cr.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}))
	defer out.Close()
	must(io.Copy(os.Stdout, out))
	log.Println("finished running", imageName)
}
