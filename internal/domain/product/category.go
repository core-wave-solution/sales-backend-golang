package productentity

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/internal/domain/entity"
)

type Category struct {
	entity.Entity
	bun.BaseModel `bun:"table:categories"`
	CategoryCommonAttributes
}

type CategoryCommonAttributes struct {
	Name                string     `bun:"name,notnull" json:"name"`
	NeedPrint           bool       `bun:"need_print,notnull" json:"need_print"`
	Sizes               []Size     `bun:"rel:has-many,join:id=category_id" json:"sizes,omitempty"`
	Quantities          []Quantity `bun:"rel:has-many,join:id=category_id" json:"quantities,omitempty"`
	Products            []Product  `bun:"rel:has-many,join:id=category_id" json:"products,omitempty"`
	Processes           []Process  `bun:"rel:has-many,join:id=category_id" json:"processes,omitempty"`
	AditionalCategories []Category `bun:"m2m:category_to_aditional_categories,join:Category=Category" json:"aditional_categories,omitempty"`
}

type PatchCategory struct {
	Name      *string `json:"name"`
	NeedPrint *bool   `json:"need_print"`
}

type CategoryToAditionalCategories struct {
	CategoryID  uuid.UUID `bun:"type:uuid,pk"`
	Category    *Category `bun:"rel:belongs-to,join:category_id=id"`
	AditionalID uuid.UUID `bun:"type:uuid,pk"`
	Aditional   *Category `bun:"rel:belongs-to,join:aditional_id=id"`
}
