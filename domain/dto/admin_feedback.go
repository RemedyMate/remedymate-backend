package dto

type ListFeedbackDTO struct {
	Limit    int    `form:"limit"`
	Offset   int    `form:"offset"`
	Language string `form:"language"`
}
