package stat

import (
	"context"
	"errors"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/dto"
)

type StatService struct {
	progressRepo domain.ProgressRepository
	goalRepo     domain.GoalRepository
}

func NewStatService(progressRepo domain.ProgressRepository, goalRepo domain.GoalRepository) *StatService {
	return &StatService{
		progressRepo: progressRepo,
		goalRepo:     goalRepo,
	}
}

func (s *StatService) GetProgress(ctx context.Context, userID, year, month int64, isMonthly bool) (dto.StatResponse, error) {
	var goalLine int64
	goal, err := s.goalRepo.GetGoalByUserID(ctx, userID)
	if errors.Is(err, domain.ErrRecordNotFound) {
		goalLine = 0
	} else if err != nil {
		return dto.StatResponse{}, err
	} else {
		goalLine = s.calculateGoalLine(goal, isMonthly)
	}

	var res []dto.Progress
	if isMonthly {
		res, err = s.progressRepo.GetMonthlyProgress(ctx, userID, year)
	} else {
		res, err = s.progressRepo.GetDailyProgress(ctx, userID, year, month)
	}
	if err != nil {
		return dto.StatResponse{}, err
	}

	return dto.StatResponse{
		Progress: res,
		Goal:     goalLine,
	}, nil
}

func (s *StatService) calculateGoalLine(goal domain.Goal, isMonthly bool) int64 {
	if goal.Type == domain.GoalTypeBooks {
		return 0
	}

	if isMonthly {
		if goal.Frequency == domain.GoalFrequencyDaily {
			return goal.Value * 30
		}
		return goal.Value
	} else {
		if goal.Frequency == domain.GoalFrequencyMonthly {
			return goal.Value / 30
		}
		return goal.Value
	}
}
