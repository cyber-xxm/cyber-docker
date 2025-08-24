package docker

import (
	"cyber-docker/internal/mods/docker/api"
	"github.com/gin-gonic/gin"
)

type Docker struct {
	ImageApi     api.Images
	ContainerApi api.Containers
	NetworkApi   api.Network
	VolumeApi    api.Volume
}

func (a *Docker) RegisterV1Routers(v1 *gin.RouterGroup) {
	v1.GET("health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	image := v1.Group("/images")
	{
		image.GET("", a.ImageApi.List)
		image.GET("/:id", a.ImageApi.Inspect)
		image.PUT("/:id", a.ImageApi.CheckUpgrade)
		image.POST("/file", a.ImageApi.Import)
		image.DELETE("", a.ImageApi.Prune)
		image.DELETE("/:id", a.ImageApi.Delete)
	}

	containers := v1.Group("/containers")
	{
		containers.GET("", a.ContainerApi.List)
		containers.GET("/:id", a.ContainerApi.Inspect)
		containers.GET("/:id/stat", a.ContainerApi.Stat)
		containers.PUT("/:id/stat", a.ContainerApi.Start)
		containers.PATCH("/:id/stat", a.ContainerApi.Stop)
		containers.GET("/:id/top", a.ContainerApi.Top)
		containers.PUT("/:id", a.ContainerApi.Update)
		containers.PUT("/:id/:name", a.ContainerApi.Commit)
		containers.GET("/:id/file", a.ContainerApi.Export)
		containers.DELETE("", a.ContainerApi.Prune)
		containers.DELETE("/:id/:name", a.ContainerApi.Delete)
	}

	networks := v1.Group("/networks")
	{
		networks.GET("", a.NetworkApi.List)
		networks.GET("/:id", a.NetworkApi.Inspect)
		networks.POST("", a.NetworkApi.Create)
		networks.PUT("/:id", a.NetworkApi.Connect)
		networks.PATCH("/:id", a.NetworkApi.Disconnect)
		networks.DELETE("", a.NetworkApi.Prune)
		networks.DELETE("/:id", a.NetworkApi.Delete)
	}

	volumes := v1.Group("/volumes")
	{
		volumes.GET("", a.VolumeApi.List)
		volumes.GET("/:id", a.VolumeApi.Inspect)
		volumes.POST("", a.VolumeApi.Create)
		volumes.DELETE("", a.VolumeApi.Prune)
		volumes.DELETE("/:id", a.VolumeApi.Delete)
	}
}
