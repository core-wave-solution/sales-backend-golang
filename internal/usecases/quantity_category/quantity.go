package quantityusecases

import (
	"context"

	"github.com/google/uuid"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	quantitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/quantity_category"
)

type Service struct {
	rq productentity.QuantityRepository
	rc productentity.CategoryRepository
}

func NewService(rq productentity.QuantityRepository, rc productentity.CategoryRepository) *Service {
	return &Service{rq: rq, rc: rc}
}

func (s *Service) RegisterQuantity(ctx context.Context, dto *quantitydto.RegisterQuantityInput) (uuid.UUID, error) {
	quantity, err := dto.ToModel()

	if err != nil {
		return uuid.Nil, err
	}

	category, err := s.rc.GetCategoryById(ctx, quantity.CategoryID.String())

	if err != nil {
		return uuid.Nil, err
	}

	if err = productentity.ValidateDuplicateQuantities(quantity.Quantity, category.Quantities); err != nil {
		return uuid.Nil, err
	}

	if err = s.rq.RegisterQuantity(ctx, quantity); err != nil {
		return uuid.Nil, err
	}

	return quantity.ID, nil
}

func (s *Service) UpdateQuantity(ctx context.Context, dtoId *entitydto.IdRequest, dto *quantitydto.UpdateQuantityInput) error {
	quantity, err := s.rq.GetQuantityById(ctx, dtoId.ID.String())

	if err != nil {
		return err
	}

	if err = dto.UpdateModel(quantity); err != nil {
		return err
	}

	category, err := s.rc.GetCategoryById(ctx, quantity.CategoryID.String())

	if err != nil {
		return err
	}

	if err = productentity.ValidateUpdateQuantity(quantity, category.Quantities); err != nil {
		return err
	}

	if err = s.rq.UpdateQuantity(ctx, quantity); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteQuantity(ctx context.Context, dto *entitydto.IdRequest) error {
	if _, err := s.rq.GetQuantityById(ctx, dto.ID.String()); err != nil {
		return err
	}

	if err := s.rq.DeleteQuantity(ctx, dto.ID.String()); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetQuantityById(ctx context.Context, dto *entitydto.IdRequest) (*productentity.Quantity, error) {
	if quantity, err := s.rq.GetQuantityById(ctx, dto.ID.String()); err != nil {
		return nil, err
	} else {
		return quantity, nil
	}
}
