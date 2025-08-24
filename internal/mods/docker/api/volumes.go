package api

import (
	"cyber-docker/internal/mods/docker/entity/dto"
	"cyber-docker/pkg/utils"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Volume struct {
	SDK *client.Client
}

func (a *Volume) List(c *gin.Context) {
	var params dto.VolumeListDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}

	filter := filters.NewArgs()
	if params.Name != "" {
		filter.Add("name", params.Name)
	}
	volumeList, err := a.SDK.VolumeList(c, volume.ListOptions{
		Filters: filter,
	})
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}

	utils.ResSuccess(c, gin.H{
		"volumeList": volumeList.Volumes,
		"warning":    volumeList.Warnings,
		"inUse":      inUse(c, a.SDK, params.Name),
	})
}

func (a *Volume) Inspect(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	var params dto.VolumeGetDto
	volumeInfo, err := a.SDK.VolumeInspect(c, id)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResSuccess(c, gin.H{
		"info":  volumeInfo,
		"inUse": inUse(c, a.SDK, params.Name),
	})
}

func (a *Volume) Create(c *gin.Context) {
	var params dto.VolumeCreateDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}
	options := make(map[string]string)
	switch params.Type {
	case "tmpfs":
		options["type"] = "tmpfs"
		options["device"] = "tmpfs"
		options["o"] = params.TmpfsOptions
	case "nfs", "nfs4":
		options["type"] = params.Type
		options["device"] = ":" + strings.TrimPrefix(params.NfsMountPoint, ":")
		options["o"] = params.NfsUrl + "," + params.NfsOptions
	case "other":
		for _, row := range params.OtherOptions {
			item := strings.Split(row, "\n")
			options[item[0]] = item[1]
		}
	}
	volumeInfo, err := a.SDK.VolumeCreate(c, volume.CreateOptions{
		Driver:     params.Driver,
		Name:       params.Name,
		DriverOpts: options,
	})
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResSuccess(c, gin.H{
		"volumeInfo": volumeInfo,
	})
}

func (a *Volume) Prune(c *gin.Context) {
	var params dto.VolumePruneDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}
	filter := filters.NewArgs()
	res, err := a.SDK.VolumesPrune(c, filter)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	// 清理非匿名未使用卷
	if params.All {
		volumeList, err := a.SDK.VolumeList(c, volume.ListOptions{})
		if err != nil {
			utils.ResError(c, http.StatusInternalServerError, err.Error())
		}
		var unUseVolume []string
		containerList, err := a.SDK.ContainerList(c, container.ListOptions{
			All:    true,
			Latest: true,
		})
		if err != nil {
			utils.ResError(c, http.StatusInternalServerError, err.Error())
		}
		for _, item := range volumeList.Volumes {
			has := false
			for _, containerInfo := range containerList {
				for _, mount := range containerInfo.Mounts {
					if mount.Name != "" && mount.Name == item.Name {
						has = true
					}
				}
			}
			if !has {
				if item.UsageData != nil {
					res.SpaceReclaimed += uint64(item.UsageData.Size)
				}
				unUseVolume = append(unUseVolume, item.Name)
			}
		}

		for _, item := range unUseVolume {
			res.VolumesDeleted = append(res.VolumesDeleted, item)
			err = a.SDK.VolumeRemove(c, item, false)
			if err != nil {
				utils.ResError(c, http.StatusInternalServerError, err.Error())
			}
		}
	}
	utils.ResSuccess(c, res)
}

func (a *Volume) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	err := a.SDK.VolumeRemove(c, id, false)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
}

func inUse(c *gin.Context, client *client.Client, name string) []map[string]interface{} {
	containerList, err := client.ContainerList(c, container.ListOptions{
		All:    true,
		Latest: true,
	})
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}

	var inUseContainer []map[string]interface{}

	for _, item := range containerList {
		for _, mount := range item.Mounts {
			if mount.Name != "" && mount.Name == name {
				inUseContainer = append(inUseContainer, map[string]interface{}{
					"name":  item.Names[0],
					"mount": mount.Destination,
					"rw":    mount.RW,
					"id":    item.ID,
				})
			}
		}
	}
	return inUseContainer

}
