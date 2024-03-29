package productentity

import (
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrSizeIsInvalid    = errors.New("size is invalid")
)

type Product struct {
	entity.Entity
	bun.BaseModel `bun:"table:products"`
	ProductCommonAttributes
}

type ProductCommonAttributes struct {
	Code        string    `bun:"code,unique,notnull" json:"code"`
	Name        string    `bun:"name,notnull" json:"name"`
	ImagePath   *string   `bun:"image_path" json:"image_path"`
	Description string    `bun:"description" json:"description"`
	Price       float64   `bun:"price,notnull" json:"price"`
	Cost        float64   `bun:"cost" json:"cost"`
	IsAvailable bool      `bun:"is_available" json:"is_available"`
	CategoryID  uuid.UUID `bun:"column:category_id,type:uuid,notnull" json:"category_id"`
	Category    *Category `bun:"rel:belongs-to" json:"category,omitempty"`
	SizeID      uuid.UUID `bun:"column:size_id,type:uuid,notnull" json:"size_id"`
	Size        *Size     `bun:"rel:belongs-to" json:"size,omitempty"`
}

type PatchProduct struct {
	Code        *string    `json:"code"`
	Name        *string    `json:"name"`
	ImagePath   *string    `json:"image_path"`
	Description *string    `json:"description"`
	Price       *float64   `json:"price"`
	Cost        *float64   `json:"cost"`
	IsAvailable *bool      `json:"is_available"`
	CategoryID  *uuid.UUID `json:"category_id"`
	SizeID      *uuid.UUID `json:"size_id"`
}

func (p *Product) FindSizeInCategory() (bool, error) {
	if p.Category == nil {
		return false, ErrCategoryNotFound
	}

	for _, v := range p.Category.Sizes {
		if v.ID == p.SizeID {
			return true, nil
		}
	}

	return false, errors.New("size not found")
}
