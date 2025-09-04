package entities

import "time"

type RefreshToken struct {
	ID        string    `bson:"_id,omitempty"`
	Token     string    `bson:"token"`
	UserID    string    `bson:"userId"`
	ExpiresAt time.Time `bson:"expiresAt"`
	CreatedAt time.Time `bson:"createdAt"`
}
