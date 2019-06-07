package docker

import (
	"context"
	"fmt"
	"strconv"

	client "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

// BasicDocker ...
type dockerType struct {
	Label       string
	Image       string
	Options     Options
	Client      *client.Client
	containerID string
}

// Options ...
type Options struct {
	ExposedPort int
	ServicePort int
}

// LocalDockerDaemon ...
type LocalDockerDaemon interface {
	Run() error
	Stop() error
	// pull()
}

// Run ...
func (d *dockerType) Run() error {
	fmt.Println("Run Docker Image")

	ctx := context.Background()

	// Pull Image first
	_, err := d.Client.ImagePull(ctx, d.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	// Create Container
	resp, err := d.Client.ContainerCreate(
		ctx,
		&container.Config{
			Image:        d.Image,
			ExposedPorts: nat.PortSet{nat.Port(strconv.Itoa(d.Options.ServicePort)): struct{}{}},
		},
		&container.HostConfig{
			PortBindings: map[nat.Port][]nat.PortBinding{
				nat.Port(strconv.Itoa(d.Options.ServicePort)): {
					{
						HostIP:   "127.0.0.1",
						HostPort: strconv.Itoa(d.Options.ExposedPort),
					},
				},
			},
		},
		nil,
		"")
	if err != nil {
		return err
	}
	d.containerID = resp.ID

	// Start Container
	if err := d.Client.ContainerStart(ctx, d.containerID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

// Stop ...
func (d dockerType) Stop() error {
	fmt.Println("Stop Docker Container")

	ctx := context.Background()

	if err := d.Client.ContainerStop(ctx, d.containerID, nil); err != nil {
		return err
	}

	if err := d.Client.ContainerRemove(ctx, d.containerID, types.ContainerRemoveOptions{RemoveVolumes: true, RemoveLinks: false, Force: false}); err != nil {
		return err
	}

	return nil
}

// NewLocalDockerDaemon ...
func NewLocalDockerDaemon(label string, image string, options Options) (LocalDockerDaemon, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	return &dockerType{label, image, options, cli, ""}, nil
}
