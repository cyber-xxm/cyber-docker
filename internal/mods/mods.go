package mods

import (
	"cyber-docker/internal/mods/docker"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

const (
	apiPrefix = "/api/"
)

type Mods struct {
	Docker *docker.Docker
}

var Set = wire.NewSet(
	wire.Struct(new(Mods), "*"),
	docker.Set,
)

func (a *Mods) RegisterRouters(e *gin.Engine) {
	gAPI := e.Group(apiPrefix)
	v1 := gAPI.Group("v1")
	a.Docker.RegisterV1Routers(v1)
}
