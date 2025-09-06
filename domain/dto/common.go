package dto

// PaginationQueryParams defines parameters for pagination.
type PaginationQueryParams struct {
	Page  int `json:"page"`  // Current page number, 1-indexed
	Limit int `json:"limit"` // Number of items per page
}

// FilterQueryParams defines parameters for filtering topics.
type FilterQueryParams struct {
	Search string `json:"search"` // General search term (e.g., for topic_key, name_en, name_am)
	Status string `json:"status"` // "active", "deleted", or "all"
}

// SortQueryParams defines parameters for sorting topics.
type SortQueryParams struct {
	SortBy string `json:"sort_by"` // Field to sort by (e.g., "topic_key", "name_en", "created_at")
	Order  string `json:"order"`   // "asc" or "desc"
}

// TopicListQueryParams combines pagination, filtering, and sorting for topic retrieval.
type TopicListQueryParams struct {
	PaginationQueryParams
	FilterQueryParams
	SortQueryParams
}

// PaginationMetadata contains metadata for paginated responses
type PaginationMetadata struct {
	Page       int   `json:"page"`        // Current page number
	Limit      int   `json:"limit"`       // Number of items per page
	Total      int64 `json:"total"`       // Total number of items
	TotalPages int   `json:"total_pages"` // Total number of pages
	HasNext    bool  `json:"has_next"`    // Whether there's a next page
	HasPrev    bool  `json:"has_prev"`    // Whether there's a previous page
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{}        `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
	Message    string             `json:"message"`
}
