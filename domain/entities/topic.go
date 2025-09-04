// domain/entities/topic.go
package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TopicStatus for soft-delete at topic level.
type TopicStatus string

const (
	TopicStatusActive  TopicStatus = "active"
	TopicStatusDeleted TopicStatus = "deleted"
)

// LocalizedGuidanceContent holds the structured guidance for one language.
type LocalizedGuidanceContent struct {
	SelfCare      []string      `json:"self_care" bson:"self_care"`
	OTCCategories []OTCCategory `json:"otc_categories,omitempty" bson:"otc_categories,omitempty"`
	SeekCareIf    []string      `json:"seek_care_if" bson:"seek_care_if"`
	Disclaimer    string        `json:"disclaimer" bson:"disclaimer"`
}

// RevisionEntry for an optional revision history (lightweight).
type RevisionEntry struct {
	Version   int                `json:"version" bson:"version"`
	Notes     string             `json:"notes,omitempty" bson:"notes,omitempty"`
	ChangedAt time.Time          `json:"changed_at" bson:"changed_at"`
	ChangedBy primitive.ObjectID `json:"changed_by,omitempty" bson:"changed_by,omitempty"`
}

// Topic is the MongoDB document for a RemedyMate topic / approved block.
type Topic struct {
	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	TopicKey      string             `json:"topic_key" bson:"topic_key"` // unique human-friendly key
	NameEN        string             `json:"name_en" bson:"name_en"`
	NameAM        string             `json:"name_am" bson:"name_am"`
	DescriptionEN string             `json:"description_en,omitempty" bson:"description_en,omitempty"`
	DescriptionAM string             `json:"description_am,omitempty" bson:"description_am,omitempty"`
	Status        TopicStatus        `json:"status" bson:"status"` // active | deleted
	Translations    map[string]LocalizedGuidanceContent `json:"translations" bson:"translations"` // expect at least "en" and "am"
	Version         int                                 `json:"version" bson:"version"`           // increment for major changes
	RevisionHistory []RevisionEntry                     `json:"revision_history,omitempty" bson:"revision_history,omitempty"`
	CreatedAt       time.Time                           `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time                           `json:"updated_at" bson:"updated_at"`
	CreatedBy       primitive.ObjectID                  `json:"created_by,omitempty" bson:"created_by,omitempty"`
	UpdatedBy       primitive.ObjectID                  `json:"updated_by,omitempty" bson:"updated_by,omitempty"`
}
