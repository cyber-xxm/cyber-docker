package dto

type ImageDto struct {
	Sha string `json:"sha" form:"sha"`
}

type ImageGetDto struct {
	Layer bool `json:"layer" form:"layer"`
}

type ImageDeleteDto struct {
	ImageDto
	Force bool `json:"force" form:"force"`
}

type ImagePruneDto struct {
	Unused bool `json:"unused" form:"unused"`
	Build  bool `json:"build" form:"build"`
}

type ImageImportDto struct {
	Container bool `json:"container,omitempty" form:"container" binding:"required"`
}
