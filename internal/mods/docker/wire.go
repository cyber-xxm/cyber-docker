package docker

import (
	"cyber-docker/internal/mods/docker/api"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Struct(new(Docker), "*"),
	wire.Struct(new(api.Images), "*"),
)
