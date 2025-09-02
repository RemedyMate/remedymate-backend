package dto

type CreateFeedbackDTO struct {
	SessionID string `json:"sessionId" binding:"required,min=3,max=128"`
	TopicKey  string `json:"topicKey" binding:"required,min=2,max=64"`
	Language  string `json:"language" binding:"required,oneof=en am"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Message   string `json:"message" binding:"omitempty,max=500"`
}
