package interfaces

import (
	"context"
	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
)

type AdminRedFlagUsecase interface {
	List(ctx context.Context) ([]entities.RedFlag, error)
	Get(ctx context.Context, id string) (*entities.RedFlag, error)
	Create(ctx context.Context, in dto.CreateRedFlagDTO, actor string) (*entities.RedFlag, error)
	Update(ctx context.Context, id string, in dto.UpdateRedFlagDTO, actor string) (*entities.RedFlag, error)
	Delete(ctx context.Context, id string, actor string) error
}

type AdminFeedbackUsecase interface {
	List(ctx context.Context, limit, offset int, language string) ([]entities.Feedback, int64, error)
	Get(ctx context.Context, id string) (*entities.Feedback, error)
	Delete(ctx context.Context, id string) error
}

type AdminAnalyticsUsecase interface {
	Get(ctx context.Context, from, to string) (map[string]interface{}, error)
}
