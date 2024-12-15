package progress

import (
	"context"
	"fmt"
	"time"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/dto"
)

type progressService struct {
	progressRepo  domain.ProgressRepository
	readingRepo   domain.ReadingRepository
	validationSvc domain.ValidationService
}

func NewProgressService(repo domain.ProgressRepository, readingRepo domain.ReadingRepository, validator domain.ValidationService) *progressService {
	return &progressService{
		progressRepo:  repo,
		readingRepo:   readingRepo,
		validationSvc: validator,
	}
}

func (s *progressService) CreateProgress(ctx context.Context, userID, readingID int64, progressReq dto.ProgressRequest) (domain.Progress, error) {
	if progressReq.Date.After(time.Now()) {
		return domain.Progress{}, fmt.Errorf("%w: %s", domain.ErrValidation, "reading date cannot be in the future")
	}

	progress := domain.Progress{
		ReadingID:   readingID,
		UserID:      userID,
		Pages:       progressReq.Pages,
		ReadingDate: progressReq.Date,
	}

	if err := s.validationSvc.ValidateStruct(progress); err != nil {
		return domain.Progress{}, err
	}

	reading, err := s.readingRepo.GetReadingByID(ctx, readingID)
	if err != nil {
		return domain.Progress{}, err
	}

	totalReadPages, err := s.progressRepo.GetTotalProgressByReadingID(ctx, readingID)
	if err != nil {
		return domain.Progress{}, err
	}

	if totalReadPages+progress.Pages > reading.TotalPages {
		return domain.Progress{}, fmt.Errorf("%w: %s", domain.ErrValidation, "total progress cannot be greater than total pages")
	}

	return s.progressRepo.CreateProgress(ctx, progress)
}
