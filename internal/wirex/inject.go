package wirex

import (
	"cyber-docker/internal/mods"
	"cyber-docker/pkg/container/di"
	"cyber-docker/pkg/docker"
	"github.com/docker/docker/client"
)

type Injector struct {
	*mods.Mods
	Client *client.Client
}

func GetDockerClient(dic *di.Container) *client.Client {
	return docker.ClientFrom(dic.Get)
}
