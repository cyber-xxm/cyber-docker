package main

import (
	"context"
	"cyber-docker/internal/bootstrap"
	"cyber-docker/pkg/container/di"
	"cyber-docker/pkg/docker"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	client, err := docker.NewDockerClientFromHost("tcp://ip:port")
	if err != nil {
		panic(err)
	}
	dic := di.NewContainer(di.ServiceConstructorMap{
		docker.ClientName: func(get di.Get) interface{} {
			return client
		},
	})
	ctx := context.Background()

	err = bootstrap.Run(ctx, dic)
	if err != nil {
		panic(err)
	}
}
