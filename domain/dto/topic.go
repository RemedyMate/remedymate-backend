package dto

import "remedymate-backend/domain/entities"

// PaginatedTopicsResult represents a paginated list of topics.
type PaginatedTopicsResult struct {
	Topics     []entities.Topic `json:"topics"`
	TotalCount int64            `json:"total_count"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
}

// TopicCreateRequest represents the data needed to create a new topic.
type TopicCreateRequest struct {
	TopicKey          string                                       `json:"topic_key"`
	NameEN            string                                       `json:"name_en"`
	NameAM            string                                       `json:"name_am"`
	DescriptionEN     string                                       `json:"description_en,omitempty"`
	DescriptionAM     string                                       `json:"description_am,omitempty"`
	IsOfflineCachable bool                                         `json:"is_offline_cachable"`
	Translations      map[string]entities.LocalizedGuidanceContent `json:"translations"` // Full content for initial creation
}

// TopicUpdateRequest represents the data for updating an existing topic.
type TopicUpdateRequest struct {
	NameEN            string                                       `json:"name_en,omitempty"`
	NameAM            string                                       `json:"name_am,omitempty"`
	DescriptionEN     string                                       `json:"description_en,omitempty"`
	DescriptionAM     string                                       `json:"description_am,omitempty"`
	IsOfflineCachable *bool                                        `json:"is_offline_cachable,omitempty"` // Pointer for explicit zero-value update
	Status            *entities.TopicStatus                        `json:"status,omitempty"`              // Pointer for explicit zero-value update
	Translations      map[string]entities.LocalizedGuidanceContent `json:"translations,omitempty"`        // Allow updating translations
}
