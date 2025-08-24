package docker

import (
	"cyber-docker/internal/mods/docker/api"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Struct(new(Docker), "*"),
	wire.Struct(new(api.Images), "*"),
	wire.Struct(new(api.Containers), "*"),
	wire.Struct(new(api.Network), "*"),
	wire.Struct(new(api.Volume), "*"),
)
