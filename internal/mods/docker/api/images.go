package api

import (
	"bufio"
	"cyber-docker/internal/mods/docker/entity/dto"
	"cyber-docker/pkg/function"
	"cyber-docker/pkg/utils"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-units"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
)

type Images struct {
	SDK *client.Client
}

func (a *Images) List(c *gin.Context) {
	imageList, err := a.SDK.ImageList(c, image.ListOptions{
		All:            true,
		ContainerCount: false,
	})
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResSuccess(c, imageList)
}

func (a *Images) Inspect(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	var params dto.ImageGetDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}

	imageDetail, err := a.SDK.ImageInspect(c, id)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	var layers []image.HistoryResponseItem
	if params.Layer {
		imageHistory, err := a.SDK.ImageHistory(c, id)
		if err != nil {
			utils.ResError(c, http.StatusInternalServerError, err.Error())
		}
		layers = append(layers, imageHistory...)
	}
	utils.ResSuccess(c, gin.H{
		"info":  imageDetail,
		"layer": layers,
	})
}

func (a *Images) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	var params dto.ImageDeleteDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}
	opts := image.RemoveOptions{
		PruneChildren: true,
	}
	if params.Force {
		opts.Force = true
	}
	_, err = a.SDK.ImageRemove(c, id, opts)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResOK(c)
}

func (a *Images) Prune(c *gin.Context) {
	var params dto.ImagePruneDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}
	if params.Build {
		res, err := a.SDK.BuildCachePrune(c, types.BuildCachePruneOptions{
			All: true,
		})
		if err != nil {
			utils.ResError(c, http.StatusInternalServerError, err.Error())
		}
		utils.ResSuccess(c, gin.H{
			"size": units.HumanSize(float64(res.SpaceReclaimed)),
		})
	}

	// 清理未使用的 tag 时，直接调用 Prune 处理
	// 只清理未使用镜像时，需要手动删除，避免 tag 被删除
	if params.Unused {
		filter := filters.NewArgs()
		filter.Add("dangling", "0")
		res, err := a.SDK.ImagesPrune(c, filter)
		if err != nil {
			utils.ResError(c, http.StatusInternalServerError, err.Error())
		}
		utils.ResSuccess(c, gin.H{
			"size":  units.HumanSize(float64(res.SpaceReclaimed)),
			"count": fmt.Sprintf("%d", len(res.ImagesDeleted)),
		})
	} else {
		var deleteImageSpaceReclaimed int64 = 0
		deleteImageTotal := 0
		useImageList := make([]string, 0)
		if containerList, err := a.SDK.ContainerList(c, container.ListOptions{}); err != nil {
			useImageList = function.PluckArrayWalk(containerList, func(item container.Summary) (string, bool) {
				return item.ImageID, true
			})
		}
		if imageList, err := a.SDK.ImageList(c, image.ListOptions{
			All:            true,
			ContainerCount: true,
		}); err != nil {
			for _, item := range imageList {
				if !function.InSlice(useImageList, item.ID) {
					deleteImageSpaceReclaimed += item.Size
					deleteImageTotal += 1
					_, _ = a.SDK.ImageRemove(c, item.ID, image.RemoveOptions{PruneChildren: true})
				}
			}
		}
		utils.ResSuccess(c, gin.H{
			"size":  units.HumanSize(float64(deleteImageSpaceReclaimed)),
			"count": fmt.Sprintf("%d", deleteImageTotal),
		})
	}
}

func (a *Images) CheckUpgrade(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}

	imageInfo, err := a.SDK.ImageInspect(c, id)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	if function.IsEmptySlice(imageInfo.RepoTags) {
		utils.ResError(c, http.StatusBadRequest, "image repo tags is empty")
	}
	// TODO 镜像仓库拿

	utils.ResOK(c)
}

func (a *Images) Import(c *gin.Context) {
	var params dto.ImageImportDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}

	file, header, err := c.Request.FormFile("file")
	fmt.Println(header.Filename)
	// 如果导入的是容器的tar包
	if params.Container {
		//var params dto.ImportDTO
		//err := c.ShouldBindBodyWithJSON(&params)
		//if err != nil {
		//	utils.ResError(c, http.StatusBadRequest, err.Error())
		//}

	} else {
		// 如果导入的是镜像的tar包
		reader := bufio.NewReader(file)
		response, err := a.SDK.ImageLoad(c, reader, client.ImageLoadWithQuiet(false))
		if err != nil {
			utils.ResError(c, http.StatusInternalServerError, err.Error())
		}
		// TODO 响应流数据
		defer func() {
			if response.Body.Close() != nil {
				slog.Error("docker", "image import ", err)
			}
		}()
		io.Copy(c.Writer, response.Body)
	}
}
