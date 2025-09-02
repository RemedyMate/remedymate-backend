package interfaces

import (
	"context"
	"remedymate-backend/domain/entities"
)

type RedFlagRepository interface {
	List(ctx context.Context) ([]entities.RedFlag, error)
	GetByID(ctx context.Context, id string) (*entities.RedFlag, error)
	Create(ctx context.Context, rf *entities.RedFlag) error
	Update(ctx context.Context, rf *entities.RedFlag) error
	SoftDelete(ctx context.Context, id string, deletedBy string) error
}

type FeedbackRepository interface {
	List(ctx context.Context, limit, offset int, language string) ([]entities.Feedback, error)
	Count(ctx context.Context, language string) (int64, error)
	GetByID(ctx context.Context, id string) (*entities.Feedback, error)
	SoftDelete(ctx context.Context, id string) error
	Create(ctx context.Context, f *entities.Feedback) error
}

type AnalyticsRepository interface {
	InsertEvent(ctx context.Context, evt map[string]interface{}) error
	UsageCounts(ctx context.Context, from, to string) (interface{}, error)
	LanguageBreakdown(ctx context.Context, from, to string) (interface{}, error)
	OfflineHits(ctx context.Context, from, to string) (interface{}, error)
	RedEscalations(ctx context.Context, from, to string) (interface{}, error)
	TopTopics(ctx context.Context, from, to string, limit int) (interface{}, error)
}
