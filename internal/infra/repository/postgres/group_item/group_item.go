package groupitemrepositorybun

import (
	"context"
	"database/sql"
	"sync"

	"github.com/uptrace/bun"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	groupitementity "github.com/willjrcom/sales-backend-go/internal/domain/group_item"
	itementity "github.com/willjrcom/sales-backend-go/internal/domain/item"
)

type GroupItemRepositoryBun struct {
	mu sync.Mutex
	db *bun.DB
}

func NewGroupItemRepositoryBun(db *bun.DB) *GroupItemRepositoryBun {
	return &GroupItemRepositoryBun{db: db}
}

func (r *GroupItemRepositoryBun) CreateGroupItem(ctx context.Context, p *groupitementity.GroupItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewInsert().Model(p).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *GroupItemRepositoryBun) UpdateGroupItem(ctx context.Context, p *groupitementity.GroupItem) error {
	p.CalculateTotalValues()

	r.mu.TryLock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	if _, err := r.db.NewUpdate().Model(p).Where("id = ?", p.ID).Exec(ctx); err != nil {
		return err
	}

	return nil
}

func (r *GroupItemRepositoryBun) CalculateTotal(ctx context.Context, id string) (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	groupItem, err := r.GetGroupByID(ctx, id, true)

	if err != nil {
		return err
	}

	groupItem.CalculateTotalValues()

	return r.UpdateGroupItem(ctx, groupItem)
}

func (r *GroupItemRepositoryBun) DeleteGroupItem(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})

	if err != nil {
		return err
	}

	if _, err = tx.NewDelete().Model(&groupitementity.GroupItem{}).Where("id = ?", id).Exec(ctx); err != nil {
		if errRoolback := tx.Rollback(); errRoolback != nil {
			return errRoolback
		}

		return err
	}

	if _, err = tx.NewDelete().Model(&itementity.Item{}).Where("group_item_id = ?", id).Exec(ctx); err != nil {
		if errRoolback := tx.Rollback(); errRoolback != nil {
			return errRoolback
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		if errRoolback := tx.Rollback(); errRoolback != nil {
			return errRoolback
		}

		return err
	}

	return nil
}

func (r *GroupItemRepositoryBun) GetGroupByID(ctx context.Context, id string, withRelation bool) (*groupitementity.GroupItem, error) {
	item := &groupitementity.GroupItem{}
	r.mu.TryLock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	query := r.db.NewSelect().Model(item).Where("group_item.id = ?", id).Relation("Category")

	if withRelation {
		query.Relation("Items")
	}

	if err := query.Scan(ctx); err != nil {
		return nil, err
	}

	return item, nil
}

func (r *GroupItemRepositoryBun) GetGroupsByStatus(ctx context.Context, status groupitementity.StatusGroupItem) ([]groupitementity.GroupItem, error) {
	items := []groupitementity.GroupItem{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&items).Where("status = ?", status).Relation("Items").Relation("Category").Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *GroupItemRepositoryBun) GetGroupsByOrderIDAndStatus(ctx context.Context, id string, status groupitementity.StatusGroupItem) ([]groupitementity.GroupItem, error) {
	items := []groupitementity.GroupItem{}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := database.ChangeSchema(ctx, r.db); err != nil {
		return nil, err
	}

	if err := r.db.NewSelect().Model(&items).Where("order_id = ? AND status = ?", id, status).Relation("Items").Relation("Category").Scan(ctx); err != nil {
		return nil, err
	}

	return items, nil
}
