package api

import (
	"cyber-docker/internal/mods/docker/entity/dto"
	"cyber-docker/pkg/function"
	"cyber-docker/pkg/utils"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Network struct {
	SDK *client.Client
}

func (a *Network) List(c *gin.Context) {
	var params dto.NetworkListDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}
	filter := filters.NewArgs()
	if params.Name != "" {
		filter.Add("name", params.Name)
	}

	networkList, err := a.SDK.NetworkList(c, network.ListOptions{
		Filters: filter,
	})
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResSuccess(c, networkList)
}

func (a *Network) Inspect(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	networkInfo, err := a.SDK.NetworkInspect(c, id, network.InspectOptions{})
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResSuccess(c, networkInfo)
}

func (a *Network) Create(c *gin.Context) {
	var params dto.NetworkCreateDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}

	option := network.CreateOptions{
		Driver: "bridge",
		Options: map[string]string{
			"name": params.Name,
		},
		EnableIPv6: function.Ptr(false),
		IPAM: &network.IPAM{
			Driver:  "default",
			Options: map[string]string{},
			Config:  []network.IPAMConfig{},
		},
	}

	if params.IpV4 != nil && params.IpV4.Gateway != "" && params.IpV4.Subnet != "" {
		option.IPAM.Config = append(option.IPAM.Config, network.IPAMConfig{
			Subnet:  params.IpV4.Subnet,
			Gateway: params.IpV4.Gateway,
		})
	}
	if params.IpV6 != nil && params.IpV6.Gateway != "" && params.IpV6.Subnet != "" {
		option.EnableIPv6 = function.Ptr(true)
		option.IPAM.Config = append(option.IPAM.Config, network.IPAMConfig{
			Subnet:  params.IpV6.Subnet,
			Gateway: params.IpV6.Gateway,
		})
	}

	response, err := a.SDK.NetworkCreate(c, params.Name, option)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResSuccess(c, response)
}

func (a *Network) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	if networkInfo, err := a.SDK.NetworkInspect(c, id, network.InspectOptions{}); err == nil {
		for _, item := range networkInfo.Containers {
			err = a.SDK.NetworkDisconnect(c, id, item.Name, true)
		}
		if err != nil {
			utils.ResError(c, http.StatusInternalServerError, err.Error())
		}
		err := a.SDK.NetworkRemove(c, id)
		if err != nil {
			utils.ResError(c, http.StatusInternalServerError, err.Error())
		}
	} else {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResOK(c)
}

func (a *Network) Connect(c *gin.Context) {
	var params dto.NetworkConnectDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}

	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}

	// 关联网络时，重新退出加入
	_ = a.SDK.NetworkDisconnect(c, id, params.ContainerName, true)
	if params.ContainerAlise == nil {
		params.ContainerAlise = make([]string, 0)
	}

	endpointSetting := &network.EndpointSettings{
		Aliases:    params.ContainerAlise,
		IPAMConfig: &network.EndpointIPAMConfig{},
		DNSNames:   make([]string, 0),
	}

	if params.IpV4 != "" {
		endpointSetting.IPAMConfig.IPv4Address = params.IpV4
	}
	if params.IpV6 != "" {
		endpointSetting.IPAMConfig.IPv6Address = params.IpV6
	}
	if !function.IsEmptyArray(params.DNSNames) {
		endpointSetting.DNSNames = params.DNSNames
	}
	if params.MacAddress != "" {
		endpointSetting.MacAddress = params.MacAddress
	}

	err = a.SDK.NetworkConnect(c, id, params.ContainerName, endpointSetting)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResOK(c)
}

func (a *Network) Disconnect(c *gin.Context) {
	var params dto.NetworkDisconnectDto
	err := c.ShouldBind(&params)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, err.Error())
	}
	id := c.Param("id")
	if id == "" {
		utils.ResError(c, http.StatusBadRequest, "id is required")
	}
	err = a.SDK.NetworkDisconnect(c, id, params.ContainerName, false)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResOK(c)
}

func (a *Network) Prune(c *gin.Context) {
	_, err := a.SDK.NetworksPrune(c, filters.NewArgs())
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, err.Error())
	}
	utils.ResOK(c)
}
