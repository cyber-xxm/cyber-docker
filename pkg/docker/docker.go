package docker

import (
	"cyber-docker/pkg/container/di"
	"github.com/docker/docker/client"
)

var ClientName = di.TypeInstanceToName(client.Client{})

func ClientFrom(get di.Get) *client.Client {
	return get(ClientName).(*client.Client)
}

func NewDockerClientFromEnv() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	return cli, err
}

func NewDockerClientFromHost(host string) (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
	return cli, err
}
