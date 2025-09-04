package entities

import "time"

type RedFlag struct {
	ID          string       `bson:"_id,omitempty" json:"id"`
	Keywords    []string     `bson:"keywords" json:"keywords"`
	Language    string       `bson:"language" json:"language"`
	Level       TriageLevel  `bson:"level" json:"level"`
	Description string       `bson:"description" json:"description"`
	IsDeleted   bool         `bson:"isDeleted" json:"-"`
	CreatedAt   time.Time    `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time    `bson:"updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time   `bson:"deletedAt,omitempty" json:"-"`
	CreatedBy   *string      `bson:"createdBy,omitempty" json:"-"`
	UpdatedBy   *string      `bson:"updatedBy,omitempty" json:"-"`
	DeletedBy   *string      `bson:"deletedBy,omitempty" json:"-"`
}
