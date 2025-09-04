package dto

type CreateRedFlagDTO struct {
	Keywords    []string `json:"keywords" binding:"required,min=1"`
	Language    string   `json:"language" binding:"required,oneof=en am"`
	Level       string   `json:"level" binding:"required,oneof=RED YELLOW"`
	Description string   `json:"description" binding:"required,min=3"`
}

type UpdateRedFlagDTO struct {
	Keywords    []string `json:"keywords" binding:"omitempty,min=1"`
	Language    string   `json:"language" binding:"omitempty,oneof=en am"`
	Level       string   `json:"level" binding:"omitempty,oneof=RED YELLOW"`
	Description string   `json:"description" binding:"omitempty,min=3"`
}
