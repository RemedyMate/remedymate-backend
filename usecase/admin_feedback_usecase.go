package usecase

import (
	"context"

	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
)

type AdminFeedbackUsecaseImpl struct {
	repo interfaces.FeedbackRepository
}

func NewAdminFeedbackUsecase(repo interfaces.FeedbackRepository) interfaces.AdminFeedbackUsecase {
	return &AdminFeedbackUsecaseImpl{repo: repo}
}

func (uc *AdminFeedbackUsecaseImpl) List(ctx context.Context, limit, offset int, language string) ([]entities.Feedback, int64, error) {
	items, err := uc.repo.List(ctx, limit, offset, language)
	if err != nil {
		return nil, 0, err
	}

	count, err := uc.repo.Count(ctx, language)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

func (uc *AdminFeedbackUsecaseImpl) Get(ctx context.Context, id string) (*entities.Feedback, error) {
	f, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (uc *AdminFeedbackUsecaseImpl) Delete(ctx context.Context, id string) error {
	return uc.repo.SoftDelete(ctx, id)
}
