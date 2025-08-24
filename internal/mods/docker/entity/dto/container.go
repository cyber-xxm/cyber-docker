package dto

type ContainerDto struct {
	Sha string `json:"sha"`
}

type ContainerUpdateDto struct {
	ContainerDto
	Name          string                  `json:"name"`
	RestartPolicy *ContainerRestartPolicy `json:"restart_policy,omitempty"`
}

type ContainerRestartPolicy struct {
	Name       string `json:"name"`
	MaxAttempt int    `json:"max_attempt"`
}

type ContainerDeleteDto struct {
	ContainerDto
	DeleteVolume bool `json:"delete_volume" binding:"required"`
	DeleteLink   bool `json:"delete_link" binding:"required"`
}

type ContainerCommitDto struct {
	ContainerDto
	Name string `json:"name" binding:"required"`
}
