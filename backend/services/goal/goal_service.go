package goal

import (
	"context"
	"errors"
	"time"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/dto"
)

type goalService struct {
	goalRepo      domain.GoalRepository
	progressRepo  domain.ProgressRepository
	readingRepo   domain.ReadingRepository
	validationSvc domain.ValidationService
}

func NewGoalService(repo domain.GoalRepository, progressRepo domain.ProgressRepository,
	readingRepo domain.ReadingRepository, validator domain.ValidationService) domain.GoalService {
	return &goalService{
		goalRepo:      repo,
		progressRepo:  progressRepo,
		readingRepo:   readingRepo,
		validationSvc: validator,
	}
}

func (s *goalService) GetGoal(ctx context.Context, userID int64) (domain.Goal, error) {
	goal, err := s.goalRepo.GetGoalByUserID(ctx, userID)
	if err != nil {
		return domain.Goal{}, err
	}

	return goal, nil
}

func (s *goalService) SetGoal(ctx context.Context, userID int64, goal domain.Goal) (domain.Goal, error) {
	currentGoal, err := s.goalRepo.GetGoalByUserID(ctx, userID)
	if err == nil {
		if goal.Type != "" {
			currentGoal.Type = goal.Type
		}
		if goal.Frequency != "" {
			currentGoal.Frequency = goal.Frequency
		}
		if goal.Value > 0 {
			currentGoal.Value = goal.Value
		}
		currentGoal.UserID = userID

		if err := s.validationSvc.ValidateStruct(currentGoal); err != nil {
			return domain.Goal{}, err
		}

		return s.goalRepo.UpdateGoal(ctx, currentGoal)
	} else if errors.Is(err, domain.ErrRecordNotFound) {
		goal.UserID = userID

		if err := s.validationSvc.ValidateStruct(goal); err != nil {
			return domain.Goal{}, err
		}

		return s.goalRepo.CreateGoal(ctx, goal)
	}

	return domain.Goal{}, err
}

func (s *goalService) GetGoalProgress(ctx context.Context, userID int64) (dto.GoalProgressResponse, error) {
	goal, err := s.goalRepo.GetGoalByUserID(ctx, userID)
	if err != nil {
		return dto.GoalProgressResponse{}, err
	}

	period := ""
	if goal.Frequency == domain.GoalFrequencyDaily {
		period = time.Now().Format("2006-01-02")
	} else if goal.Frequency == domain.GoalFrequencyMonthly {
		period = time.Now().Format("2006-01")
	}
	readingIDs, err := s.progressRepo.GetUserReadingIDsByPeriod(ctx, userID, period)

	if err != nil {
		return dto.GoalProgressResponse{}, err
	}

	var progress int64
	for _, readingID := range readingIDs {
		if goal.Type == domain.GoalTypePages {
			dayProgress, err := s.progressRepo.GetProgressByReadingAndDate(ctx, readingID, period)
			if err != nil {
				return dto.GoalProgressResponse{}, err
			}

			progress += dayProgress
		} else if goal.Type == domain.GoalTypeBooks {
			totalProgress, err := s.progressRepo.GetTotalProgressByReadingID(ctx, readingID)
			if err != nil {
				return dto.GoalProgressResponse{}, err
			}

			reading, err := s.readingRepo.GetReadingByID(ctx, readingID)
			if err != nil {
				return dto.GoalProgressResponse{}, err
			}

			if reading.GetStatus(totalProgress) == domain.ReadingStatusCompleted {
				progress++
			}
		}
	}

	progress = min(progress, goal.Value)
	goalProgress := dto.GoalProgressResponse{
		Percentage: float64(progress) / float64(goal.Value) * 100,
		Left:       goal.Value - progress,
	}

	return goalProgress, nil
}
