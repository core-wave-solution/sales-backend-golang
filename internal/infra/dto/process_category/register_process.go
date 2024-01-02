package processdto

import (
	"errors"

	"github.com/google/uuid"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
)

var (
	ErrNameRequired              = errors.New("name is required")
	ErrOrderRequired             = errors.New("order is required")
	ErrIdealTimeRequired         = errors.New("ideal time is required")
	ErrExperimentalErrorRequired = errors.New("experimental error is required")
	ErrCategoryRequired          = errors.New("category ID is required")
)

type RegisterProcessInput struct {
	productentity.ProcessCommonAttributes
}

func (s *RegisterProcessInput) validate() error {
	if s.Name == "" {
		return ErrNameRequired
	}
	if s.Order < 1 {
		return ErrOrderRequired
	}

	if s.IdealTime == nil {
		return ErrIdealTimeRequired
	}

	if s.ExperimentalError == nil {
		return ErrExperimentalErrorRequired
	}

	if s.CategoryID == uuid.Nil {
		return ErrCategoryRequired
	}

	return nil
}

func (s *RegisterProcessInput) ToModel() (*productentity.Process, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}
	processCommonAttributes := productentity.ProcessCommonAttributes{
		Name:              s.Name,
		Order:             s.Order,
		IdealTime:         s.IdealTime,
		ExperimentalError: s.ExperimentalError,
		CategoryID:        s.CategoryID,
	}

	return &productentity.Process{
		Entity:                  entity.NewEntity(),
		ProcessCommonAttributes: processCommonAttributes,
	}, nil
}