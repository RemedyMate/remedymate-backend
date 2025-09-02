package entities

import "time"

type Feedback struct {
	ID        string     `bson:"_id,omitempty" json:"id"`
	SessionID string     `bson:"sessionId" json:"sessionId"`
	TopicKey  string     `bson:"topicKey" json:"topicKey"`
	Language  string     `bson:"language" json:"language"`
	Rating    int        `bson:"rating" json:"rating"`
	Message   string     `bson:"message" json:"message"`
	IsDeleted bool       `bson:"isDeleted" json:"-"`
	CreatedAt time.Time  `bson:"createdAt" json:"createdAt"`
	DeletedAt *time.Time `bson:"deletedAt,omitempty" json:"-"`
}
