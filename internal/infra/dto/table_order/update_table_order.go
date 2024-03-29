package tableorderdto

import (
	"github.com/google/uuid"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
)

type UpdateTableOrderInput struct {
	orderentity.TableOrderCommonAttributes
}

func (o *UpdateTableOrderInput) validate() error {
	return nil
}

func (s *UpdateTableOrderInput) ToModel() (model *orderentity.TableOrder, err error) {
	if err = s.validate(); err != nil {
		return nil, err
	}

	tableCommonAttributes := orderentity.TableOrderCommonAttributes{}

	if s.WaiterID != uuid.Nil {
		model.WaiterID = s.WaiterID
	}

	if s.TableID != uuid.Nil {
		model.TableID = s.TableID
	}

	return &orderentity.TableOrder{
		TableOrderCommonAttributes: tableCommonAttributes,
	}, err
}
