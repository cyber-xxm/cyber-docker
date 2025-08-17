package docker

import (
	"cyber-docker/internal/mods/docker/api"
	"github.com/gin-gonic/gin"
)

type Docker struct {
	ImageApi api.Images
}

func (a *Docker) RegisterV1Routers(v1 *gin.RouterGroup) {
	v1.GET("health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	image := v1.Group("images")
	{
		image.GET("", a.ImageApi.List)
		//image.GET(":id", a.ImageApi.Get)
		//image.GET("pages", a.ImageApi.GetAllPage)
		//image.POST("", a.ImageApi.Create)
		//image.PUT(":id", a.ImageApi.Update)
		//image.DELETE(":id", a.ImageApi.Delete)
	}
}
