package addressentity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Address struct {
	entity.Entity
	bun.BaseModel `bun:"table:addresses"`
	AddressCommonAttributes
}

type AddressCommonAttributes struct {
	ObjectID     uuid.UUID `bun:"object_id,type:uuid,notnull" json:"object_id"`
	Street       string    `bun:"street,notnull" json:"street"`
	Number       string    `bun:"number,notnull" json:"number"`
	Complement   string    `bun:"complement" json:"complement"`
	Reference    string    `bun:"reference" json:"reference"`
	Neighborhood string    `bun:"neighborhood,notnull" json:"neighborhood"`
	City         string    `bun:"city,notnull" json:"city"`
	State        string    `bun:"state,notnull" json:"state"`
	Cep          string    `bun:"cep" json:"cep"`
	DeliveryTax  float64   `bun:"delivery_tax,notnull" json:"delivery_tax"`
}

type PatchAddress struct {
	Street       *string  `json:"street"`
	Number       *string  `json:"number"`
	Complement   *string  `json:"complement"`
	Reference    *string  `json:"reference"`
	Neighborhood *string  `json:"neighborhood"`
	City         *string  `json:"city"`
	State        *string  `json:"state"`
	Cep          *string  `json:"cep"`
	DeliveryTax  *float64 `json:"delivery_tax"`
}

func (a *Address) Validate() error {
	if a.Street == "" {
		return errors.New("street is required")
	}
	if a.Number == "" {
		return errors.New("number is required")
	}
	if a.Neighborhood == "" {
		return errors.New("neighborhood is required")
	}
	if a.City == "" {
		return errors.New("city is required")
	}
	if a.State == "" {
		return errors.New("state is required")
	}
	return nil
}

func NewAddress(addressCommonAttributes *AddressCommonAttributes) *Address {
	return &Address{
		Entity:                  entity.NewEntity(),
		AddressCommonAttributes: *addressCommonAttributes,
	}
}
