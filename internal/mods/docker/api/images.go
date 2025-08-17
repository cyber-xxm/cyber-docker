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
	"log/slog"
	"net/http"
)

type Images struct {
	Client *client.Client
}

func (i *Images) List(ctx *gin.Context) {
	context := ctx.Request.Context()
	imageList, err := i.Client.ImageList(context, image.ListOptions{
		All:            true,
		ContainerCount: false,
	})
	if err != nil {
		utils.ResError(ctx, http.StatusInternalServerError, err.Error())
	}
	utils.ResSuccess(ctx, imageList)
}

func (i *Images) Get(ctx *gin.Context, sha string, layer bool) {
	context := ctx.Request.Context()
	imageDetail, err := i.Client.ImageInspect(context, sha)
	if err != nil {
		utils.ResError(ctx, http.StatusInternalServerError, err.Error())
	}
	var layers []image.HistoryResponseItem
	if layer {
		imageHistory, err := i.Client.ImageHistory(context, sha)
		if err != nil {
			utils.ResError(ctx, http.StatusInternalServerError, err.Error())
		}
		layers = append(layers, imageHistory...)
	}
	utils.ResSuccess(ctx, gin.H{
		"info":  imageDetail,
		"layer": layers,
	})
}

func (i *Images) Delete(ctx *gin.Context, sha string, force bool) {
	context := ctx.Request.Context()
	opts := image.RemoveOptions{
		PruneChildren: true,
	}
	if force {
		opts.Force = true
	}
	_, err := i.Client.ImageRemove(context, sha, opts)
	if err != nil {
		utils.ResError(ctx, http.StatusInternalServerError, err.Error())
	}
	utils.ResOK(ctx)
}

func (i *Images) Prune(ctx *gin.Context, unused, build bool) {
	context := ctx.Request.Context()
	if build {
		res, err := i.Client.BuildCachePrune(context, types.BuildCachePruneOptions{
			All: true,
		})
		if err != nil {
			utils.ResError(ctx, http.StatusInternalServerError, err.Error())
		}
		utils.ResSuccess(ctx, gin.H{
			"size": units.HumanSize(float64(res.SpaceReclaimed)),
		})
	}

	// 清理未使用的 tag 时，直接调用 Prune 处理
	// 只清理未使用镜像时，需要手动删除，避免 tag 被删除
	if unused {
		filter := filters.NewArgs()
		filter.Add("dangling", "0")
		res, err := i.Client.ImagesPrune(context, filter)
		if err != nil {
			utils.ResError(ctx, http.StatusInternalServerError, err.Error())
		}
		utils.ResSuccess(ctx, gin.H{
			"size":  units.HumanSize(float64(res.SpaceReclaimed)),
			"count": fmt.Sprintf("%d", len(res.ImagesDeleted)),
		})
	} else {
		var deleteImageSpaceReclaimed int64 = 0
		deleteImageTotal := 0
		useImageList := make([]string, 0)
		if containerList, err := i.Client.ContainerList(context, container.ListOptions{}); err != nil {
			useImageList = function.PluckArrayWalk(containerList, func(item container.Summary) (string, bool) {
				return item.ImageID, true
			})
		}
		if imageList, err := i.Client.ImageList(context, image.ListOptions{
			All:            true,
			ContainerCount: true,
		}); err != nil {
			for _, item := range imageList {
				if !function.InSlice(useImageList, item.ID) {
					deleteImageSpaceReclaimed += item.Size
					deleteImageTotal += 1
					_, _ = i.Client.ImageRemove(context, item.ID, image.RemoveOptions{PruneChildren: true})
				}
			}
		}
		utils.ResSuccess(ctx, gin.H{
			"size":  units.HumanSize(float64(deleteImageSpaceReclaimed)),
			"count": fmt.Sprintf("%d", deleteImageTotal),
		})
	}
}

func (i *Images) CheckUpgrade(ctx *gin.Context, sha string) {
	context := ctx.Request.Context()
	imageInfo, err := i.Client.ImageInspect(context, sha)
	if err != nil {
		utils.ResError(ctx, http.StatusInternalServerError, err.Error())
	}
	if function.IsEmptySlice(imageInfo.RepoTags) {
		utils.ResError(ctx, http.StatusBadRequest, "image repo tags is empty")
	}
	// TODO 镜像仓库拿

	utils.ResOK(ctx)
}

func (i *Images) Import(ctx *gin.Context, container bool) {
	context := ctx.Request.Context()
	file, header, err := ctx.Request.FormFile("file")
	fmt.Println(header.Filename)
	if err != nil {
		utils.ResError(ctx, http.StatusInternalServerError, err.Error())
	}
	// 如果导入的是容器的tar包
	if container {
		var params dto.ImportDTO
		err := ctx.ShouldBindBodyWithJSON(&params)
		if err != nil {
			utils.ResError(ctx, http.StatusBadRequest, err.Error())
		}

	} else {
		// 如果导入的是镜像的tar包
		reader := bufio.NewReader(file)
		response, err := i.Client.ImageLoad(context, reader, client.ImageLoadWithQuiet(false))
		if err != nil {
			utils.ResError(ctx, http.StatusInternalServerError, err.Error())
		}
		// TODO 响应流数据
		defer func() {
			if response.Body.Close() != nil {
				slog.Error("docker", "image import ", err)
			}
		}()
	}
}
