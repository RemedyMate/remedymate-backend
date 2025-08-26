package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/RemedyMate/remedymate-backend/domain/dto"
	"github.com/RemedyMate/remedymate-backend/domain/entities"
	"github.com/RemedyMate/remedymate-backend/domain/interfaces"
)

type RemedyMateUsecase struct {
	triageService  interfaces.TriageService
	contentService interfaces.ContentService
}

//  creates a new RemedyMate usecase
func NewRemedyMateUsecase(
	triageService interfaces.TriageService,
	contentService interfaces.ContentService,

) interfaces.RemedyMateUsecase {
	return &RemedyMateUsecase{
		triageService:  triageService,
		contentService: contentService,
	}
}

// performs only triage classification
func (rmu *RemedyMateUsecase) GetTriage(ctx context.Context, req dto.TriageRequest) (*dto.TriageResponse, error) {
	input := entities.SymptomInput{
		Text:     req.Text,
		Language: req.Language,
	}

	result, err := rmu.triageService.ClassifySymptoms(ctx, input)
	if err != nil {
		return nil, err
	}

	return &dto.TriageResponse{
		Level:     result.Level,
		RedFlags:  result.RedFlags,
		Message:   result.Message,
		SessionID: generateSessionID(),
	}, nil
}

// generates a unique session ID
func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}
