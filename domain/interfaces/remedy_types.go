package interfaces

// RemedyResponse represents the remedy response structure
type RemedyResponse struct {
	SessionID string         `json:"session_id"`
	Triage    TriageResponse `json:"triage"`
	Content   *GuidanceCard  `json:"guidance_card,omitempty"`
}

// TriageResponse represents triage result
type TriageResponse struct {
	Level     string   `json:"level"`
	RedFlags  []string `json:"red_flags"`
	Message   string   `json:"message"`
	SessionID string   `json:"session_id,omitempty"`
}

// GuidanceCard represents guidance card
type GuidanceCard struct {
	TopicKey      string   `json:"topic_key"`
	Language      string   `json:"language"`
	SelfCare      []string `json:"self_care"`
	OTCCategories []string `json:"otc_categories"`
	SeekCareIf    []string `json:"seek_care_if"`
	Disclaimer    string   `json:"disclaimer"`
	IsOffline     bool     `json:"is_offline"`
}
