package dto

type ImportDTO struct {
	Tar      string   `json:"tar" binding:"required"`
	Tag      string   `json:"tag" binding:"required"`
	Registry string   `json:"registry"`
	Cmd      string   `json:"cmd" binding:"required"`
	WorkDir  string   `json:"workDir"`
	Expose   []string `json:"expose"`
	Env      []string `json:"env"`
	Volume   []string `json:"volume"`
}
