package usecase

import (
	"context"
	"strings"

	"remedymate-backend/domain/dto"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
)

type AdminRedFlagUsecaseImpl struct {
	repo interfaces.RedFlagRepository
}

func NewAdminRedFlagUsecase(repo interfaces.RedFlagRepository) interfaces.AdminRedFlagUsecase {
	return &AdminRedFlagUsecaseImpl{repo: repo}
}

func (uc *AdminRedFlagUsecaseImpl) List(ctx context.Context) ([]entities.RedFlag, error) {
	redFlags, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return redFlags, nil
}

func (uc *AdminRedFlagUsecaseImpl) Get(ctx context.Context, id string) (*entities.RedFlag, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *AdminRedFlagUsecaseImpl) Create(ctx context.Context, in dto.CreateRedFlagDTO, actor string) (*entities.RedFlag, error) {
	level := entities.TriageLevel(strings.ToUpper(in.Level))
	rf := &entities.RedFlag{
		Keywords:    in.Keywords,
		Language:    in.Language,
		Level:       level,
		Description: in.Description,
		CreatedBy:   &actor,
	}
	if err := uc.repo.Create(ctx, rf); err != nil {
		return nil, err
	}
	return rf, nil
}

func (uc *AdminRedFlagUsecaseImpl) Update(ctx context.Context, id string, in dto.UpdateRedFlagDTO, actor string) (*entities.RedFlag, error) {
	existing, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Keywords != nil {
		existing.Keywords = in.Keywords
	}
	if in.Language != "" {
		existing.Language = in.Language
	}
	if in.Level != "" {
		existing.Level = entities.TriageLevel(strings.ToUpper(in.Level))
	}
	if in.Description != "" {
		existing.Description = in.Description
	}
	existing.UpdatedBy = &actor
	if err := uc.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (uc *AdminRedFlagUsecaseImpl) Delete(ctx context.Context, id string, actor string) error {
	if actor == "" {
		actor = "system"
	}
	if err := uc.repo.SoftDelete(ctx, id, actor); err != nil {
		return err
	}
	return nil
}
