package dto

type VolumeDto struct {
	Sha string `json:"sha"`
}

type VolumeListDto struct {
	VolumeDto
	Name string `json:"name" form:"name"`
}

type VolumeGetDto struct {
	VolumeDto
	Name string `json:"name" form:"name"`
}

type VolumeCreateDto struct {
	Name          string   `json:"name" binding:"required"`
	Driver        string   `json:"driver" binding:"omitempty,oneof=local"`
	Type          string   `json:"type"`
	TmpfsOptions  string   `json:"tmpfs_options"`
	NfsMountPoint string   `json:"nfsMountPoint"`
	NfsUrl        string   `json:"nfsUrl"`
	NfsOptions    string   `json:"nfsOptions"`
	OtherOptions  []string `json:"otherOptions"`
}

type VolumePruneDto struct {
	VolumeDto
	All bool `json:"all" form:"all"`
}
