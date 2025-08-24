package api

import (
	"bufio"
	"cyber-docker/internal/mods/docker/entity/dto"
	"cyber-docker/pkg/utils"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
)

var (
	restartPolicyMap = map[string]container.RestartPolicyMode{
		"always":         container.RestartPolicyAlways,
		"no":             container.RestartPolicyDisabled,
		"unless-stopped": container.RestartPolicyUnlessStopped,
		"on-failure":     container.RestartPolicyOnFailure,
	}
)

type Containers struct {
	SDK *client.Client
}

func (a *Containers) List(c *gin.Context) {
	var params dto.ContainerDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}

	filter := filters.NewArgs()
	if params.Sha != "" {
		filter.Add("id", params.Sha)
	}

	imageList, err := a.SDK.ContainerList(c, container.ListOptions{
		All:     true,
		Latest:  true,
		Filters: filter,
	})
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResSuccess(c, imageList)
}

func (a *Containers) Inspect(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	detail, err := a.SDK.ContainerInspect(c, id)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResSuccess(c, detail)
}

func (a *Containers) Start(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	err := a.SDK.ContainerStart(c, id, container.StartOptions{})
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResOK(c)
}

func (a *Containers) Stop(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	err := a.SDK.ContainerStop(c, id, container.StopOptions{})
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResOK(c)
}

func (a *Containers) Stat(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}

	response, err := a.SDK.ContainerStats(c, id, true)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	io.Copy(c.Writer, response.Body)
}

func (a *Containers) Top(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	response, err := a.SDK.ContainerTop(c, id, nil)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResSuccess(c, response)
}

func (a *Containers) Update(c *gin.Context) {
	var params dto.ContainerUpdateDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	if params.RestartPolicy != nil {
		restartPolicy := container.RestartPolicy{}
		if params.RestartPolicy.Name != "" {
			if name, ok := restartPolicyMap[params.RestartPolicy.Name]; ok {
				restartPolicy.Name = name
			} else {
				restartPolicy.Name = container.RestartPolicyDisabled
			}
		}
		if restartPolicy.Name == container.RestartPolicyOnFailure {
			restartPolicy.MaximumRetryCount = 5
		}
		if params.RestartPolicy.MaxAttempt > 0 {
			restartPolicy.Name = container.RestartPolicyOnFailure
			restartPolicy.MaximumRetryCount = params.RestartPolicy.MaxAttempt
		}

		_, err := a.SDK.ContainerUpdate(c, id, container.UpdateConfig{
			RestartPolicy: restartPolicy,
		})
		if err != nil {
			utils.ResError(c, http.StatusInternalServerError, err.Error())
		}
	}
	if params.Name != "" {
		err = a.SDK.ContainerRename(c, id, params.Name)
		if err != nil {
			utils.ResError(c, http.StatusInternalServerError, err.Error())
		}
	}
	utils.ResOK(c)
}

func (a *Containers) Prune(c *gin.Context) {
	info, err := a.SDK.ContainersPrune(c, filters.NewArgs())
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResSuccess(c, info)
}

func (a *Containers) Delete(c *gin.Context) {
	var params dto.ContainerDeleteDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	containerInfo, err := a.SDK.ContainerInspect(c, id)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}

	err = a.SDK.ContainerStop(c, id, container.StopOptions{})
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}

	err = a.SDK.ContainerRemove(c, id, container.RemoveOptions{
		RemoveVolumes: params.DeleteVolume,
		RemoveLinks:   params.DeleteLink,
	})

	if params.DeleteVolume {
		for _, item := range containerInfo.Mounts {
			if item.Type == mount.TypeVolume {
				err = a.SDK.VolumeRemove(c, item.Name, false)
				if err != nil {
					slog.Debug("remove container volume", err.Error())
				}
			}
		}
	}
	utils.ResOK(c)
}

func (a *Containers) Export(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}

	out, err := a.SDK.ContainerExport(c, id)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	defer func() {
		_ = out.Close()
	}()

	reader := bufio.NewReader(out)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	c.Header("Content-Type", "application/tar")
	c.Header("Content-Disposition", "attachment; filename="+id+".tar")
	io.Copy(c.Writer, reader)
}

func (a *Containers) Commit(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	name := c.Param("name")
	out, err := a.SDK.ContainerCommit(c, id, container.CommitOptions{
		Reference: name,
	})
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}

	utils.ResSuccess(c, out)
}
