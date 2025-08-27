package dto

type MapTopicRequest struct {
	UserInput string `json:"user_input" binding:"required"`
	Language  string `json:"language"`
}

type MapTopicResponse struct {
	TopicKey string `json:"topic_key"`
}
