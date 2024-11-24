package reading

import (
	"context"
	"fmt"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/dto"
	"github.com/rimvydascivilis/book-tracker/backend/utils"
)

type ReadingService struct {
	readingRepo   domain.ReadingRepository
	progressRepo  domain.ProgressRepository
	bookRepo      domain.BookRepository
	validationSvc domain.ValidationService
}

func NewReadingService(repo domain.ReadingRepository, progressRepo domain.ProgressRepository,
	bookRepo domain.BookRepository, validationSvc domain.ValidationService) *ReadingService {
	return &ReadingService{
		readingRepo:   repo,
		progressRepo:  progressRepo,
		bookRepo:      bookRepo,
		validationSvc: validationSvc,
	}
}

func (s *ReadingService) GetReadings(ctx context.Context, userID, page, limit int64) ([]dto.ReadingResponse, bool, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	totalCount, err := s.readingRepo.CountReadingsByUserID(ctx, userID)
	if err != nil {
		return nil, false, err
	}

	hasMore := totalCount > page*limit
	if totalCount == 0 {
		return []dto.ReadingResponse{}, false, nil
	}

	offset := (page - 1) * limit
	readings, err := s.readingRepo.GetReadingsByUserID(ctx, userID, offset, limit)
	if err != nil {
		return nil, false, err
	}

	combinedResponse := make([]dto.ReadingResponse, 0, len(readings))
	for _, reading := range readings {
		book, err := s.bookRepo.GetBookByUserID(ctx, userID, reading.BookID)
		if err != nil {
			return nil, false, err
		}

		progress, err := s.progressRepo.GetTotalProgressByReadingID(ctx, reading.ID)
		if err != nil {
			return nil, false, err
		}

		combinedResponse = append(combinedResponse, dto.ReadingResponse{
			BookTitle: book.Title,
			Status:    reading.GetStatus(progress),
			Progress:  progress,
			Reading: dto.Reading{
				ID:         reading.ID,
				TotalPages: reading.TotalPages,
				Link:       reading.Link,
			},
		})
	}

	return combinedResponse, hasMore, nil
}

func (s *ReadingService) CreateReading(ctx context.Context, userID int64, reading domain.Reading) (domain.Reading, error) {
	reading.UserID = userID
	reading.CreatedAt = utils.Now()
	reading.UpdatedAt = utils.Now()
	if err := s.validationSvc.ValidateStruct(reading); err != nil {
		return domain.Reading{}, err
	}

	_, err := s.bookRepo.GetBookByUserID(ctx, userID, reading.BookID)
	if err != nil {
		return domain.Reading{}, err
	}

	readingCount, err := s.readingRepo.CountReadingsByUserIDAndBookID(ctx, userID, reading.BookID)
	if err != nil {
		return domain.Reading{}, err
	}
	if readingCount > 0 {
		return domain.Reading{}, fmt.Errorf("%w: %d", domain.ErrAlreadyExists, reading.BookID)
	}

	return s.readingRepo.CreateReading(ctx, reading)
}
