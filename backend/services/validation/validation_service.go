package validation

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type validationService struct {
	validator *validator.Validate
}

func NewValidationService() *validationService {
	return &validationService{
		validator: validator.New(),
	}
}

func (v *validationService) ValidateStruct(s interface{}) error {
	if err := v.validator.Struct(s); err != nil {
		return v.formatValidationError(err)
	}
	return nil
}

func (v *validationService) formatValidationError(err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var errorMessages []string
		for _, fe := range ve {
			errorMessages = append(errorMessages, formatFieldError(fe))
		}
		return fmt.Errorf("%w: %s", domain.ErrValidation, errorMessages)
	}
	return domain.ErrValidation
}

func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required.", fe.Field())
	case "max":
		return fmt.Sprintf("The %s field must be at most %s characters long.", fe.Field(), fe.Param())
	case "min":
		return fmt.Sprintf("The %s field must be at least %s characters long.", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("The %s field is invalid.", fe.Field())
	}
}
