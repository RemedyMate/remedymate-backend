package entities

// RedFlagRule represents a rule for detecting red flag symptoms
type RedFlagRule struct {
	Keywords    []string    `json:"keywords" bson:"keywords"`
	Language    string      `json:"language" bson:"language"`
	Level       TriageLevel `json:"level" bson:"level"`
	Description string      `json:"description" bson:"description"`
}

// GuidanceCard represents the final guidance card shown to users
type GuidanceCard struct {
	TopicKey      string        `json:"topic_key" bson:"topic_key"`
	Language      string        `json:"language" bson:"language"`
	SelfCare      []string      `json:"self_care" bson:"self_care"`
	OTCCategories []OTCCategory `json:"otc_categories,omitempty" bson:"otc_categories,omitempty"`
	SeekCareIf    []string      `json:"seek_care_if" bson:"seek_care_if"`
	Disclaimer    string        `json:"disclaimer" bson:"disclaimer"`
	IsOffline     bool          `json:"is_offline" bson:"is_offline"`
}

// OTCCategory represents an over-the-counter medication category
type OTCCategory struct {
	CategoryName string `json:"category_name" bson:"category_name"`
	SafetyNote   string `json:"safety_note" bson:"safety_note"`
}

// ContentTranslation represents content in a specific language
type ContentTranslation struct {
	SelfCare      []string      `json:"self_care" bson:"self_care"`
	OTCCategories []OTCCategory `json:"otc_categories" bson:"otc_categories"`
	SeekCareIf    []string      `json:"seek_care_if" bson:"seek_care_if"`
	Disclaimer    string        `json:"disclaimer" bson:"disclaimer"`
}

// ApprovedBlock represents a topic with its content in multiple languages
type ApprovedBlock struct {
	TopicKey     string                        `json:"topic_key" bson:"topic_key"`
	Translations map[string]ContentTranslation `json:"translations" bson:"translations"`
}
