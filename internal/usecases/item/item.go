package itemusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
	orderentity "github.com/willjrcom/sales-backend-go/internal/domain/order"
	productentity "github.com/willjrcom/sales-backend-go/internal/domain/product"
	entitydto "github.com/willjrcom/sales-backend-go/internal/infra/dto/entity"
	itemdto "github.com/willjrcom/sales-backend-go/internal/infra/dto/item"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrSizeNotFound     = errors.New("size not found")
	ErrGroupNotStaging  = errors.New("group not staging")
)

type Service struct {
	ri  itementity.ItemRepository
	rgi groupitementity.GroupItemRepository
	ro  orderentity.OrderRepository
	rp  productentity.ProductRepository
	rq  productentity.QuantityRepository
}

func NewService(ri itementity.ItemRepository, rgi groupitementity.GroupItemRepository, ro orderentity.OrderRepository, rp productentity.ProductRepository, rq productentity.QuantityRepository) *Service {
	return &Service{ri: ri, rgi: rgi, ro: ro, rp: rp, rq: rq}
}

func (s *Service) AddItemOrder(ctx context.Context, dto *itemdto.AddItemOrderInput) (ids *itemdto.ItemIDAndGroupItemOutput, err error) {
	if _, err := s.ro.GetOrderById(ctx, dto.OrderID.String()); err != nil {
		return nil, err
	}

	product, err := s.rp.GetProductById(ctx, dto.ProductID.String())

	if err != nil {
		return nil, err
	}

	if product.Category == nil {
		return nil, ErrCategoryNotFound
	}

	if product.Size == nil {
		return nil, ErrSizeNotFound
	}

	if dto.GroupItemID == nil {
		groupItem, err := s.newGroupItem(ctx, dto.OrderID, product)

		if err != nil {
			return nil, err
		}

		dto.GroupItemID = &groupItem.ID
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, dto.GroupItemID.String(), false)

	if err != nil {
		return nil, err
	}

	if !groupItem.CanAddItems() {
		return nil, ErrGroupNotStaging
	}

	quantity, err := s.rq.GetQuantityById(ctx, dto.QuantityID.String())

	if err != nil {
		return nil, err
	}

	item, err := dto.ToModel(product, groupItem, quantity)

	if err != nil {
		return nil, err
	}

	if err = s.ri.AddItem(ctx, item); err != nil {
		return nil, err
	}

	if err = s.rgi.CalculateTotal(ctx, groupItem.ID.String()); err != nil {
		return nil, err
	}

	return &itemdto.ItemIDAndGroupItemOutput{
		GroupItemID: groupItem.ID,
		ItemID:      item.ID,
	}, nil
}

func (s *Service) StartItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = item.StartItem(); err != nil {
		return err
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), false)

	if err != nil {
		return err
	}

	if groupItem.Status == groupitementity.StatusGroupPending {
		groupItem.StartGroupItem()
	}

	if err = s.ri.UpdateItem(ctx, item); err != nil {
		return err
	}

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}

func (s *Service) ReadyItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = item.ReadyItem(); err != nil {
		return err
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), false)

	if err != nil {
		return err
	}

	if err = s.ri.UpdateItem(ctx, item); err != nil {
		return err
	}

	isAllItemsReady := true
	for i := range groupItem.Items {
		if groupItem.Items[i].Status != itementity.StatusItemReady {
			isAllItemsReady = false
			break
		}
	}

	if !isAllItemsReady {
		return nil
	}

	if err = groupItem.ReadyGroupItem(); err != nil {
		return err
	}

	return s.rgi.UpdateGroupItem(ctx, groupItem)
}

func (s *Service) CancelItem(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	item.CancelItem()

	return s.ri.UpdateItem(ctx, item)
}

func (s *Service) DeleteItemOrder(ctx context.Context, dto *entitydto.IdRequest) (err error) {
	item, err := s.ri.GetItemById(ctx, dto.ID.String())

	if err != nil {
		return err
	}

	if err = s.ri.DeleteItem(ctx, dto.ID.String()); err != nil {
		return err
	}

	groupItem, err := s.rgi.GetGroupByID(ctx, item.GroupItemID.String(), true)

	if err != nil {
		return err
	}

	if len(groupItem.Items) == 0 {
		if err = s.rgi.DeleteGroupItem(ctx, item.GroupItemID.String()); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) AddAditionalItemOrder(ctx context.Context, dto *entitydto.IdRequest) (id uuid.UUID, err error) {
	return uuid.Nil, nil
}

func (s *Service) newGroupItem(ctx context.Context, orderID uuid.UUID, product *productentity.Product) (groupItem *groupitementity.GroupItem, err error) {
	groupCommonAttributes := groupitementity.GroupCommonAttributes{
		OrderID: orderID,
		GroupDetails: groupitementity.GroupDetails{
			CategoryID: product.CategoryID,
			Size:       product.Size.Name,
			NeedPrint:  product.Category.NeedPrint,
		},
	}

	groupItem = groupitementity.NewGroupItem(groupCommonAttributes)
	err = s.rgi.CreateGroupItem(ctx, groupItem)
	return
}
