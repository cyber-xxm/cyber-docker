package bootstrap

import (
	"context"
	"cyber-docker/internal/wirex"
	"cyber-docker/pkg/container/di"
)

func Run(ctx context.Context, dic *di.Container) error {
	injector, err := wirex.BuildInjector(dic)
	if err != nil {
		panic(err)
	}
	err = startHTTPServer(injector)
	if err != nil {
		return err
	}
	return nil
}
