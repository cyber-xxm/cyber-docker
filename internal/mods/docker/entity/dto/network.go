package dto

type NetworkDto struct {
	Sha string `json:"sha"`
}

type NetworkListDto struct {
	NetworkDto
	Name string `json:"name" form:"name"`
}

type NetworkCreateDto struct {
	Name string             `json:"name"`
	IpV4 *NetworkCreateItem `json:"ipV4"`
	IpV6 *NetworkCreateItem `json:"IpV6"`
}

type NetworkCreateItem struct {
	Address string `json:"address"`
	Subnet  string `json:"subnet"`
	Gateway string `json:"gateway"`
}

type NetworkConnectDto struct {
	NetworkDto
	Name           string   `json:"name" binding:"required"`
	ContainerName  string   `json:"container_name" binding:"required"`
	ContainerAlise []string `json:"containerAlise"`
	IpV4           string   `json:"ipV4"`
	IpV6           string   `json:"ipV6"`
	MacAddress     string   `json:"mac_address"`
	DNSNames       []string `json:"dnsNames"`
}

type NetworkDisconnectDto struct {
	NetworkDto
	ContainerName string `json:"container_name" binding:"required"`
}
