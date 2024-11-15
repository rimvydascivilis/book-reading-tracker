package goal

import (
	"context"
	"errors"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type goalService struct {
	goalRepo      domain.GoalRepository
	validationSvc domain.ValidationService
}

func NewGoalService(repo domain.GoalRepository, validator domain.ValidationService) domain.GoalService {
	return &goalService{
		goalRepo:      repo,
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

		if err := s.validationSvc.ValidateStruct(currentGoal); err != nil {
			return domain.Goal{}, err
		}

		return s.goalRepo.UpdateGoal(ctx, userID, currentGoal)
	} else if errors.Is(err, domain.ErrRecordNotFound) {
		goal.UserID = userID

		if err := s.validationSvc.ValidateStruct(goal); err != nil {
			return domain.Goal{}, err
		}

		return s.goalRepo.CreateGoal(ctx, userID, goal)
	}

	return domain.Goal{}, err
}
