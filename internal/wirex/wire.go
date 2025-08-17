//go:build wireinject
// +build wireinject

package wirex

// The build tag makes sure the stub is not built in the final build.

import (
	"cyber-docker/internal/mods"
	"cyber-docker/pkg/container/di"

	"github.com/google/wire"
)

func BuildInjector(dic *di.Container) (*Injector, error) {
	wire.Build(
		GetDockerClient,
		wire.NewSet(wire.Struct(new(Injector), "*")),
		mods.Set,
	) // end
	return new(Injector), nil
}
