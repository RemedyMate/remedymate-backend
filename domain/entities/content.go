package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

// TranslationCategory represents each OTC category inside translations
type TranslationCategory struct {
	CategoryName string `bson:"category_name" json:"category_name"`
	SafetyNote   string `bson:"safety_note" json:"safety_note"`
}

// LanguageTranslation represents translation for a single language
type LanguageTranslation struct {
	SelfCare      []string              `bson:"self_care" json:"self_care"`
	OtcCategories []TranslationCategory `bson:"otc_categories" json:"otc_categories"`
	SeekCareIf    []string              `bson:"seek_care_if" json:"seek_care_if"`
	Disclaimer    string                `bson:"disclaimer" json:"disclaimer"`
}

// Translations holds both English and Amharic translations
type Translations struct {
	En LanguageTranslation `bson:"en" json:"en"`
	Am LanguageTranslation `bson:"am" json:"am"`
}

// HeadacheEntity (or HealthTopic) represents the full MongoDB document
type HealthTopic struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	TopicKey      string             `bson:"topic_key" json:"topic_key"`
	NameEn        string             `bson:"name_en" json:"name_en"`
	NameAm        string             `bson:"name_am" json:"name_am"`
	DescriptionEn string             `bson:"description_en" json:"description_en"`
	DescriptionAm string             `bson:"description_am" json:"description_am"`
	Status        string             `bson:"status" json:"status"`
	Translations  Translations       `bson:"translations" json:"translations"`
	Version       int                `bson:"version" json:"version"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy     primitive.ObjectID `bson:"created_by" json:"created_by"`
	UpdatedBy     primitive.ObjectID `bson:"updated_by" json:"updated_by"`
}
