package pickuporderusecases

import (
	"context"

	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
)

func (s *Service) LaunchOrder(ctx context.Context, dtoID *entitydto.IdRequest) (err error) {
	pickupOrder, err := s.rp.GetPickupById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err := pickupOrder.Launch(); err != nil {
		return err
	}

	if err = s.rp.UpdatePickupOrder(ctx, pickupOrder); err != nil {
		return err
	}

	return nil
}

func (s *Service) PickupOrder(ctx context.Context, dtoID *entitydto.IdRequest) (err error) {
	pickupOrder, err := s.rp.GetPickupById(ctx, dtoID.ID.String())

	if err != nil {
		return err
	}

	if err := pickupOrder.PickUp(); err != nil {
		return err
	}

	if err = s.rp.UpdatePickupOrder(ctx, pickupOrder); err != nil {
		return err
	}

	return nil
}
